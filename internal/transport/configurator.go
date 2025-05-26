package transport

import (
	"context"
	"errors"

	"github.com/DanilNaum/gophermarket/internal/api/models"
	"github.com/DanilNaum/gophermarket/internal/api/restapi/operations"
	"github.com/DanilNaum/gophermarket/internal/api/restapi/operations/balance"
	"github.com/DanilNaum/gophermarket/internal/api/restapi/operations/orders"
	"github.com/DanilNaum/gophermarket/internal/api/restapi/operations/user"
	uc "github.com/DanilNaum/gophermarket/internal/usecase"
	ucmodel "github.com/DanilNaum/gophermarket/internal/usecase/model"
	"github.com/DanilNaum/gophermarket/pkg/luhn"
	openapierrors "github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type usecase interface {
	ParseToken(token string) (uuid.UUID, error)
	Register(ctx context.Context, login, password string) (string, error)
	Login(ctx context.Context, login, password string) (string, error)

	AddOrder(ctx context.Context, userID uuid.UUID, orderID string) error
	GetOrders(ctx context.Context, userID uuid.UUID) ([]*ucmodel.Order, error)

	Balance(ctx context.Context, userID uuid.UUID) (current int, withdrawn int, err error)
	Withdraw(ctx context.Context, userID uuid.UUID, orderID string, accrual int) error
	GetUserWithdrawals(ctx context.Context, userID uuid.UUID) ([]*ucmodel.Withdrawal, error)
}

type configurator struct {
	usecase usecase
	log     *zap.SugaredLogger
}

func NewConfigurator(usecase usecase, log *zap.SugaredLogger) (*configurator, error) {
	if usecase == nil {
		return nil, errors.New("usecase is nil")
	}

	return &configurator{usecase: usecase, log: log}, nil
}

func (c *configurator) Configure(api *operations.GopherMarketLoyaltySystemAPIAPI) error {
	api.BearerAuthAuth = c.auth
	api.CookieAuthAuth = c.auth

	api.UserPostAPIUserRegisterHandler = user.PostAPIUserRegisterHandlerFunc(c.register)
	api.UserPostAPIUserLoginHandler = user.PostAPIUserLoginHandlerFunc(c.login)

	api.OrdersPostAPIUserOrdersHandler = orders.PostAPIUserOrdersHandlerFunc(c.ordersPostAPIUserOrdersHandler)
	api.OrdersGetAPIUserOrdersHandler = orders.GetAPIUserOrdersHandlerFunc(c.ordersGetAPIUserOrdersHandler)

	api.BalanceGetAPIUserWithdrawalsHandler = balance.GetAPIUserWithdrawalsHandlerFunc(c.withdrawHistory)
	api.BalanceGetAPIUserBalanceHandler = balance.GetAPIUserBalanceHandlerFunc(c.balance)
	api.BalancePostAPIUserBalanceWithdrawHandler = balance.PostAPIUserBalanceWithdrawHandlerFunc(c.withdraw)

	return nil
}

func (c *configurator) auth(token string) (interface{}, error) {
	uuid, err := c.usecase.ParseToken(token)
	if err != nil {
		c.log.Debug(err.Error())
		return nil, openapierrors.Unauthenticated("")

	}
	return uuid, nil
}

func (c *configurator) register(params user.PostAPIUserRegisterParams) middleware.Responder {
	token, err := c.usecase.Register(params.HTTPRequest.Context(), *params.Body.Login, *params.Body.Password)

	if err != nil {
		switch {
		case errors.Is(err, uc.ErrConflict):
			return user.NewPostAPIUserRegisterConflict()
		default:
			return user.NewPostAPIUserRegisterInternalServerError()
		}
	}
	return user.NewPostAPIUserRegisterOK().WithAuthorization(token).WithSetCookie(token)
}

func (c *configurator) login(params user.PostAPIUserLoginParams) middleware.Responder {
	token, err := c.usecase.Login(params.HTTPRequest.Context(), *params.Body.Login, *params.Body.Password)
	if err != nil {
		switch {
		case errors.Is(err, uc.ErrNotFound) || errors.Is(err, uc.ErrInvalidPassword):
			return user.NewPostAPIUserLoginUnauthorized()
		default:
			return user.NewPostAPIUserLoginInternalServerError()
		}
	}
	return user.NewPostAPIUserLoginOK().WithAuthorization(token).WithSetCookie(token)
}

func (c *configurator) ordersPostAPIUserOrdersHandler(params orders.PostAPIUserOrdersParams, principal interface{}) middleware.Responder {
	userID, err := getUserID(principal)
	if err != nil {
		return orders.NewPostAPIUserOrdersInternalServerError()
	}

	orderID := params.Body

	if !luhn.LuhnCheck(orderID) {
		return orders.NewPostAPIUserOrdersUnprocessableEntity()
	}

	err = c.usecase.AddOrder(params.HTTPRequest.Context(), userID, orderID)
	if err != nil {
		switch {
		case errors.Is(err, uc.ErrExists):
			return orders.NewPostAPIUserOrdersOK()
		case errors.Is(err, uc.ErrConflict):
			return orders.NewPostAPIUserOrdersConflict()
		default:
			return orders.NewPostAPIUserOrdersInternalServerError()
		}
	}

	return orders.NewPostAPIUserOrdersAccepted()
}

func (c *configurator) ordersGetAPIUserOrdersHandler(param orders.GetAPIUserOrdersParams, principal interface{}) middleware.Responder {
	userID, err := getUserID(principal)
	if err != nil {
		return orders.NewPostAPIUserOrdersInternalServerError()
	}

	odrs, err := c.usecase.GetOrders(param.HTTPRequest.Context(), userID)
	if err != nil {
		return orders.NewPostAPIUserOrdersInternalServerError()
	}

	if len(odrs) == 0 {
		return orders.NewGetAPIUserOrdersNoContent()
	}

	return orders.NewGetAPIUserOrdersOK().WithPayload(ordersFromUCModel(odrs))
}

func (c *configurator) withdrawHistory(param balance.GetAPIUserWithdrawalsParams, principal interface{}) middleware.Responder {
	userID, err := getUserID(principal)
	if err != nil {
		return balance.NewGetAPIUserWithdrawalsInternalServerError()
	}
	ws, err := c.usecase.GetUserWithdrawals(param.HTTPRequest.Context(), userID)
	if err != nil {
		return balance.NewGetAPIUserWithdrawalsInternalServerError()
	}
	if len(ws) == 0 {
		return balance.NewGetAPIUserWithdrawalsNoContent()
	}

	return balance.NewGetAPIUserWithdrawalsOK().WithPayload(withdrawsFromUCModel(ws))

}

func (c *configurator) balance(param balance.GetAPIUserBalanceParams, principal interface{}) middleware.Responder {
	userID, err := getUserID(principal)
	if err != nil {
		return balance.NewGetAPIUserBalanceInternalServerError()
	}
	current, withdraw, err := c.usecase.Balance(param.HTTPRequest.Context(), userID)
	if err != nil {
		return balance.NewGetAPIUserBalanceInternalServerError()
	}

	return balance.NewGetAPIUserBalanceOK().WithPayload(&models.Balance{
		Current:   float64(current),
		Withdrawn: float64(withdraw),
	})
}

func (c *configurator) withdraw(param balance.PostAPIUserBalanceWithdrawParams, principal interface{}) middleware.Responder {
	userID, err := getUserID(principal)
	if err != nil {
		return balance.NewPostAPIUserBalanceWithdrawInternalServerError()
	}

	orderID := *param.Body.Order

	if !luhn.LuhnCheck(orderID) {
		return balance.NewPostAPIUserBalanceWithdrawUnprocessableEntity()
	}

	err = c.usecase.Withdraw(param.HTTPRequest.Context(), userID, *param.Body.Order, int(*param.Body.Sum))
	if err != nil {
		switch {
		case errors.Is(err, uc.ErrNotEnough):
			return balance.NewPostAPIUserBalanceWithdrawPaymentRequired()
		}
		return balance.NewPostAPIUserBalanceWithdrawInternalServerError()
	}
	return balance.NewPostAPIUserBalanceWithdrawOK()
}

func getUserID(principal interface{}) (uuid.UUID, error) {
	userID, ok := principal.(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("invalid principal type")
	}
	return userID, nil
}

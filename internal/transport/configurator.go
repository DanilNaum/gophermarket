package transport

import (
	"context"
	"errors"

	"time"

	"github.com/DanilNaum/gophermarket/internal/api/models"
	"github.com/DanilNaum/gophermarket/internal/api/restapi/operations"
	"github.com/DanilNaum/gophermarket/internal/api/restapi/operations/orders"
	"github.com/DanilNaum/gophermarket/internal/api/restapi/operations/user"
	uc "github.com/DanilNaum/gophermarket/internal/usecase"
	openapierrors "github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
)

type usecase interface {
	ParseToken(token string) (uuid.UUID, error)
	Register(ctx context.Context, login, password string) (string, error)
	Login(ctx context.Context, login, password string) (string, error)
}

type configurator struct {
	usecase usecase
}

func NewConfigurator(usecase usecase) (*configurator, error) {
	if usecase == nil {
		return nil, errors.New("usecase is nil")
	}

	return &configurator{usecase: usecase}, nil
}

func (c *configurator) Configure(api *operations.GopherMarketLoyaltySystemAPIAPI) error {
	api.BearerAuthAuth = c.auth
	api.CookieAuthAuth = c.auth

	api.UserPostAPIUserRegisterHandler = user.PostAPIUserRegisterHandlerFunc(c.register)
	api.UserPostAPIUserLoginHandler = user.PostAPIUserLoginHandlerFunc(c.login)

	api.OrdersGetAPIUserOrdersHandler = orders.GetAPIUserOrdersHandlerFunc(c.ordersGetAPIUserOrdersHandler)
	return nil
}

func (c *configurator) auth(token string) (interface{}, error) {
	uuid, err := c.usecase.ParseToken(token)
	if err != nil {
		// todo: подумать насколько нужно отдавать текст ошибки
		return nil, openapierrors.Unauthenticated(err.Error())

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

func (c *configurator) ordersGetAPIUserOrdersHandler(param orders.GetAPIUserOrdersParams, principal interface{}) middleware.Responder {
	ords := []*models.Order{{
		Accrual:    500,
		Number:     "123",
		Status:     "NEW",
		UploadedAt: strfmt.DateTime(time.Now()),
	}}
	return orders.NewGetAPIUserOrdersOK().WithPayload(ords)
}

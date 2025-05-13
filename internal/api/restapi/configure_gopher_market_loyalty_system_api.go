// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target ..\..\api --name GopherMarketLoyaltySystemAPI --spec ..\..\..\swagger.json --principal interface{} --exclude-main

package restapi

// Library say to configure endpoints here,
// but I want to configure them in new package
// so code is not used

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/DanilNaum/gophermarket/internal/api/restapi/operations"
	"github.com/DanilNaum/gophermarket/internal/api/restapi/operations/balance"
	"github.com/DanilNaum/gophermarket/internal/api/restapi/operations/orders"
	"github.com/DanilNaum/gophermarket/internal/api/restapi/operations/user"
)

func configureFlags(api *operations.GopherMarketLoyaltySystemAPIAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.GopherMarketLoyaltySystemAPIAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()
	api.TxtConsumer = runtime.TextConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "Authorization" header is set
	if api.BearerAuthAuth == nil {
		api.BearerAuthAuth = func(token string) (interface{}, error) {
			return nil, errors.NotImplemented("api key auth (bearerAuth) Authorization from header param [Authorization] has not yet been implemented")
		}
	}
	// Applies when the "Cookie" header is set
	if api.CookieAuthAuth == nil {
		api.CookieAuthAuth = func(token string) (interface{}, error) {
			return nil, errors.NotImplemented("api key auth (cookieAuth) Cookie from header param [Cookie] has not yet been implemented")
		}
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()

	if api.BalanceGetAPIUserBalanceHandler == nil {
		api.BalanceGetAPIUserBalanceHandler = balance.GetAPIUserBalanceHandlerFunc(func(params balance.GetAPIUserBalanceParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation balance.GetAPIUserBalance has not yet been implemented")
		})
	}
	if api.OrdersGetAPIUserOrdersHandler == nil {
		api.OrdersGetAPIUserOrdersHandler = orders.GetAPIUserOrdersHandlerFunc(func(params orders.GetAPIUserOrdersParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation orders.GetAPIUserOrders has not yet been implemented")
		})
	}
	if api.BalanceGetAPIUserWithdrawalsHandler == nil {
		api.BalanceGetAPIUserWithdrawalsHandler = balance.GetAPIUserWithdrawalsHandlerFunc(func(params balance.GetAPIUserWithdrawalsParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation balance.GetAPIUserWithdrawals has not yet been implemented")
		})
	}
	if api.BalancePostAPIUserBalanceWithdrawHandler == nil {
		api.BalancePostAPIUserBalanceWithdrawHandler = balance.PostAPIUserBalanceWithdrawHandlerFunc(func(params balance.PostAPIUserBalanceWithdrawParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation balance.PostAPIUserBalanceWithdraw has not yet been implemented")
		})
	}
	if api.UserPostAPIUserLoginHandler == nil {
		api.UserPostAPIUserLoginHandler = user.PostAPIUserLoginHandlerFunc(func(params user.PostAPIUserLoginParams) middleware.Responder {
			return middleware.NotImplemented("operation user.PostAPIUserLogin has not yet been implemented")
		})
	}
	if api.OrdersPostAPIUserOrdersHandler == nil {
		api.OrdersPostAPIUserOrdersHandler = orders.PostAPIUserOrdersHandlerFunc(func(params orders.PostAPIUserOrdersParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation orders.PostAPIUserOrders has not yet been implemented")
		})
	}
	if api.UserPostAPIUserRegisterHandler == nil {
		api.UserPostAPIUserRegisterHandler = user.PostAPIUserRegisterHandlerFunc(func(params user.PostAPIUserRegisterParams) middleware.Responder {
			return middleware.NotImplemented("operation user.PostAPIUserRegister has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}

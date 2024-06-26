package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// Apply requirePermission middleware to 5 movies routes
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.requirePermission("movies:read", app.listMoviesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.requirePermission("movies:write", app.createMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.requirePermission("movies:read", app.showMovieHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.requirePermission("movies:write", app.updateMovieHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.requirePermission("movies:write", app.deleteMovieHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	// Add the login and logout endpoints
	// router.HandlerFunc(http.MethodPost, "/v1/login", app.loginUserHandler)
	// router.HandlerFunc(http.MethodPost, "/v1/logout", app.logoutUserHandler)

	router.HandlerFunc(http.MethodPut, "/v1/users/password", app.updateUserPasswordHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/password-reset", app.createPasswordResetTokenHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/activation", app.createActivationTokenHandler)

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	// Wrap the router
	// 1. Middleware recovery panic
	// 2. Middleware rate limit for http request
	// 3. Middleware authenticate for http request
	// 4. Middleware enableCORS for http request
	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}

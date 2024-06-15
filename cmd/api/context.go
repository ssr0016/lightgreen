package main

import (
	"context"
	"net/http"

	"greenlight.samson.net/internal/data"
)

type contextKey string

// Assign contextKey type for "user" to userContextKey constant.
const userContextKey = contextKey("user")

// contextSetUser adds User struct to request context with userContextKey.
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// contextGetUser retrieves User struct from request context.
func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}

package chttp

import (
	"net/http"

	"github.com/tusharsoni/copper"
)

// MiddlewareFunc can be used to create a middleware that can be used on a route.
type MiddlewareFunc func(http.Handler) http.Handler

// GlobalMiddlewareFuncResult can be provided to the application container to register middlewares to all routes.
type GlobalMiddlewareFuncResult struct {
	copper.Module

	GlobalMiddlewareFunc MiddlewareFunc `group:"chttp/global_middlewares"`
}

func NewGlobalMiddleware(mw MiddlewareFunc) GlobalMiddlewareFuncResult {
	return GlobalMiddlewareFuncResult{
		GlobalMiddlewareFunc: mw,
	}
}

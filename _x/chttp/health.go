package chttp

import (
	"net/http"

	"go.uber.org/fx"
)

// HealthRouteParams holds the dependencies needed to create the Health route using NewHealthRoute.
type HealthRouteParams struct {
	fx.In

	Config Config `optional:"true"`
}

// NewHealthRoute provides a route that responds OK to signify the health of the server.
func NewHealthRoute(p HealthRouteParams) RouteResult {
	config := p.Config
	if !config.isValid() {
		config = GetDefaultConfig()
	}

	route := Route{
		MiddlewareFuncs: []MiddlewareFunc{},
		Path:            config.HealthPath,
		Methods:         []string{http.MethodGet},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("OK"))
		}),
	}
	return RouteResult{Route: route}
}

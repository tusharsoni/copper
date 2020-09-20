package chttp

import (
	"net/http"
)

func NewHealthRouter(config Config) Router {
	ro := healthRouter{}

	return NewRouter([]Route{
		{
			MiddlewareFuncs: []MiddlewareFunc{},
			Path:            config.HealthPath,
			Methods:         []string{http.MethodGet},
			Handler:         http.HandlerFunc(ro.Handle),
		},
	})
}

type healthRouter struct{}

func (ro *healthRouter) Handle(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("OK"))
}

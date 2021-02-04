package main

import (
	"net/http"

	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
)

func NewAppRouter(rw chttp.ReaderWriter, logger clogger.Logger) chttp.Router {
	ro := AppRouter{
		rw:     rw,
		logger: logger,
	}

	return chttp.NewRouter([]chttp.Route{
		{
			Path:    "/",
			Methods: []string{http.MethodGet},
			Handler: http.HandlerFunc(ro.HandleHelloWorld),
		},
	})
}

type AppRouter struct {
	rw     chttp.ReaderWriter
	logger clogger.Logger
}

func (ro *AppRouter) HandleHelloWorld(w http.ResponseWriter, r *http.Request) {
	ro.rw.OK(w, map[string]string{
		"response": "Hello, World!",
	})
}

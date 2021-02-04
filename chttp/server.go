package chttp

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/tusharsoni/copper"

	"github.com/gorilla/mux"
	"github.com/tusharsoni/copper/clogger"
)

type StartServerParams struct {
	copper.StartFuncParams

	Routes            []Route          `group:"chttp/routes"`
	GlobalMiddlewares []MiddlewareFunc `group:"chttp/global_middlewares"`
	Config            Config           `optional:"true"`
	Lifecycle         copper.Lifecycle
	Logger            clogger.Logger
}

func StartServer(p StartServerParams) {
	var (
		muxRouter  = mux.NewRouter()
		muxHandler = http.NewServeMux()
		config     = p.Config
	)

	if !config.isValid() {
		config = GetDefaultConfig()
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: muxHandler,
	}

	for _, f := range p.GlobalMiddlewares {
		muxRouter.Use(mux.MiddlewareFunc(f))
	}

	if len(p.Routes) == 0 {
		p.Logger.Warn("No routes to register", nil)
	}

	sortRoutes(p.Routes)

	for _, route := range p.Routes {
		p.Logger.WithTags(map[string]interface{}{
			"path":    route.Path,
			"methods": strings.Join(route.Methods, ", "),
		}).Info("Registering route..")

		handlerFunc := route.Handler

		for _, f := range route.MiddlewareFuncs {
			handlerFunc = f(handlerFunc)
		}

		muxRoute := muxRouter.Handle(route.Path, handlerFunc)

		if len(route.Methods) > 0 {
			muxRoute.Methods(route.Methods...)
		}
	}

	muxHandler.Handle("/", muxRouter)

	go func() {
		p.Logger.WithTags(map[string]interface{}{
			"port": config.Port,
		}).Info("Starting http server..")

		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			p.Logger.Error("Failed to start http server", err)
		}
	}()

	p.Lifecycle.OnStop(func(ctx context.Context) error {
		return httpServer.Shutdown(ctx)
	})
}

package chttp

import (
	"github.com/tusharsoni/copper"
	"github.com/tusharsoni/copper/clogger"
)

type Module struct {
	copper.Module

	RequestLoggerMiddleware GlobalMiddlewareFuncResult
	JSONReaderWriter        ReaderWriter
	HealthRouter            Router
}

type NewParams struct {
	copper.ModuleParams

	Config Config `optional:"true"`
	Logger clogger.Logger
}

func New(p NewParams) Module {
	config := p.Config
	if !config.isValid() {
		config = GetDefaultConfig()
	}

	return Module{
		RequestLoggerMiddleware: GlobalMiddlewareFuncResult{
			GlobalMiddlewareFunc: NewRequestLogger(p.Logger),
		},
		JSONReaderWriter: NewJSONReaderWriter(p.Logger),
		HealthRouter:     NewHealthRouter(config),
	}
}

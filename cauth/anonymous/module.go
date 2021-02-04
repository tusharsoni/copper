package anonymous

import (
	"github.com/tusharsoni/copper"
	"github.com/tusharsoni/copper/cauth"
	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
)

type Module struct {
	copper.Module

	Svc    Svc
	Router chttp.Router
}

type NewParams struct {
	copper.ModuleParams

	Auth   cauth.Svc
	RW     chttp.ReaderWriter
	Logger clogger.Logger
}

func New(p NewParams) Module {
	svc := NewSvc(SvcParams{
		Auth: p.Auth,
	})

	return Module{
		Svc: svc,
		Router: NewRouter(RouterParams{
			Svc:    svc,
			RW:     p.RW,
			Logger: p.Logger,
		}),
	}
}

package cauth

import (
	"github.com/tusharsoni/copper"
	"github.com/tusharsoni/copper/cacl"
	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
	"gorm.io/gorm"
)

type Module struct {
	copper.Module

	Svc               Svc
	SessionMiddleware SessionMiddleware
}

type NewParams struct {
	copper.ModuleParams

	DB     *gorm.DB
	RW     chttp.ReaderWriter
	Logger clogger.Logger

	ACL cacl.Svc `optional:"true"`
}

func New(p NewParams) Module {
	var (
		svc = NewSvc(NewSQLRepo(p.DB))
		mw  = NewSessionMiddleware(NewSessionMiddlewareParams{
			RW:     p.RW,
			Svc:    svc,
			Logger: p.Logger,
			ACL:    p.ACL,
		})
	)

	return Module{
		Svc:               svc,
		SessionMiddleware: mw,
	}
}

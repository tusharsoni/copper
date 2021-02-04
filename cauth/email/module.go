package email

import (
	"github.com/tusharsoni/copper"
	"github.com/tusharsoni/copper/cauth"
	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
	"github.com/tusharsoni/copper/cmailer"
	"gorm.io/gorm"
)

type Module struct {
	copper.Module

	Svc    Svc
	Router chttp.Router
}

type NewParams struct {
	copper.ModuleParams

	Auth      cauth.Svc
	Mailer    cmailer.Mailer
	Logger    clogger.Logger
	Config    Config
	RW        chttp.ReaderWriter
	SessionMW cauth.SessionMiddleware

	DB *gorm.DB
}

func New(p NewParams) Module {
	svc := NewSvc(SvcParams{
		Auth:   p.Auth,
		Repo:   NewSQLRepo(p.DB),
		Mailer: p.Mailer,
		Config: p.Config,
		Logger: p.Logger,
	})

	return Module{
		Svc: svc,
		Router: NewRouter(NewRouterParams{
			RW:        p.RW,
			Logger:    p.Logger,
			Auth:      svc,
			SessionMW: chttp.MiddlewareFunc(p.SessionMW),
			Config:    p.Config,
		}),
	}
}

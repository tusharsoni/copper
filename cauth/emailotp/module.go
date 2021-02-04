package emailotp

import (
	"github.com/tusharsoni/copper"
	"github.com/tusharsoni/copper/cauth"
	"github.com/tusharsoni/copper/cerror"
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

	Auth   cauth.Svc
	RW     chttp.ReaderWriter
	Mailer cmailer.Mailer
	Logger clogger.Logger
	Config Config

	DB *gorm.DB
}

func New(p NewParams) (Module, error) {
	svc, err := NewSvc(NewSvcParams{
		Auth:   p.Auth,
		Repo:   NewSQLRepo(p.DB),
		Mailer: p.Mailer,
		Config: p.Config,
	})
	if err != nil {
		return Module{}, cerror.New(err, "failed to create service", nil)
	}

	return Module{
		Svc: svc,
		Router: NewRouter(NewRouterParams{
			RW:     p.RW,
			Logger: p.Logger,
			Auth:   svc,
		}),
	}, nil
}

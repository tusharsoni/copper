package cqueue

import (
	"github.com/tusharsoni/copper"
	"github.com/tusharsoni/copper/chttp"

	"github.com/tusharsoni/copper/clogger"
	"gorm.io/gorm"
)

type Module struct {
	copper.Module

	Svc    Svc
	Router chttp.Router
}

type NewParams struct {
	copper.ModuleParams

	DB     *gorm.DB
	RW     chttp.ReaderWriter
	Logger clogger.Logger
	Config Config `optional:"true"`
}

func New(p NewParams) Module {
	config := p.Config
	if !config.isValid() {
		config = GetDefaultConfig()
	}

	svc := NewSvc(SvcParams{
		Repo:   NewSQLRepo(p.DB),
		Config: config,
	})

	return Module{
		Svc: svc,
		Router: NewRouter(RouterParams{
			RW:     p.RW,
			Logger: p.Logger,
			Svc:    svc,
		}),
	}
}

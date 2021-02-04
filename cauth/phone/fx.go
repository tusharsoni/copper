package phone

import (
	"context"
	"github.com/tusharsoni/copper"
	"github.com/tusharsoni/copper/cauth"
	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
	"github.com/tusharsoni/copper/ctexter"
	"gorm.io/gorm"
)

type Module struct {
	copper.Module

	Svc    Svc
	Repo   Repo
	Router chttp.Router
}

type NewParams struct {
	copper.ModuleParams

	DB *gorm.DB
	Auth cauth.Svc
	Texter ctexter.Svc
	RW chttp.ReaderWriter
	Logger clogger.Logger

	LC copper.Lifecycle
}

func New(p NewParams) Module {
	var (
		repo = NewSQLRepo(p.DB)
		svc = NewSvc(p.Auth, repo, p.Texter)
	)

	p.LC.OnStart(func(ctx context.Context) error {
		AddPhoneNumberValidator(p.Logger)
		return nil
	})

	return Module{
		Svc:    NewSvc(p.Auth, repo, p.Texter),
		Repo:   repo,
		Router: NewRouter(RouterParams{
			RW:     p.RW,
			Logger: p.Logger,
			Auth:   svc,
		}),
	}
}

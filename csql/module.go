package csql

import (
	"github.com/tusharsoni/copper"
	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
	"gorm.io/gorm"
)

type Module struct {
	copper.Module

	DB            *gorm.DB
	TxnMiddleware chttp.GlobalMiddlewareFuncResult
}

type NewParams struct {
	copper.ModuleParams

	Config Config
	Logger clogger.Logger
}

func New(p NewParams) (Module, error) {
	db, err := NewGormDB(GormDBParams{
		Config: p.Config,
		Logger: p.Logger,
	})
	if err != nil {
		return Module{}, cerror.New(err, "failed to create database connection", nil)
	}

	return Module{
		DB:            db,
		TxnMiddleware: chttp.NewGlobalMiddleware(NewTxnMiddleware(db, p.Logger)),
	}, nil
}

// Package copper provides the primitives to create a new app using github.com/uber-go/fx.
package copper

import (
	"context"
	"log"

	"github.com/tusharsoni/copper/crandom"
	"go.uber.org/fx"
)

type ModuleParams fx.In
type Module fx.Out

type App struct {
	opts []fx.Option
}

func New() *App {
	return &App{
		opts: []fx.Option{
			fx.Invoke(crandom.Seed),
		},
	}
}

func (a *App) AddModules(constructors ...interface{}) {
	for i := range constructors {
		a.opts = append(a.opts, fx.Provide(constructors[i]))
	}
}

func (a *App) AddConfigs(configs ...interface{}) {
	for i := range configs {
		a.opts = append(a.opts, fx.Supply(configs[i]))
	}
}

func (a *App) Start(funcs ...interface{}) {
	for i := range funcs {
		a.opts = append(a.opts, fx.Invoke(funcs[i]))
	}

	fx.New(a.opts...).Run()
}

func (a *App) Run(funcs ...interface{}) {
	for i := range funcs {
		a.opts = append(a.opts, fx.Invoke(funcs[i]))
	}

	fxApp := fx.New(a.opts...)

	startCtx, cancel := context.WithTimeout(context.Background(), fxApp.StartTimeout())
	defer cancel()

	if err := fxApp.Start(startCtx); err != nil {
		log.Fatalf("ERROR\t\tFailed to start: %v", err)
	}

	stopCtx, cancel := context.WithTimeout(context.Background(), fxApp.StopTimeout())
	defer cancel()

	if err := fxApp.Stop(stopCtx); err != nil {
		log.Fatalf("ERROR\t\tFailed to stop cleanly: %v", err)
	}
}

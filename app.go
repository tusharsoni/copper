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

type StartFuncParams fx.In

type App struct {
	opts []fx.Option
}

func New() *App {
	return &App{
		opts: []fx.Option{
			fx.Provide(newLifecycle),
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

// Start starts the application and blocks until it receives a kill signal. Then, it shuts down gracefully.
func (a *App) Start(funcs ...interface{}) {
	for i := range funcs {
		a.opts = append(a.opts, fx.Invoke(funcs[i]))
	}

	fx.New(a.opts...).Run()
}

// Run runs all of the provided funcs in the app container and shuts down gracefully.
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

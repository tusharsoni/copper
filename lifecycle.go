package copper

import (
	"context"

	"go.uber.org/fx"
)

type Lifecycle interface {
	OnStart(f func(context.Context) error)
	OnStop(f func(context.Context) error)
}

func newLifecycle(lc fx.Lifecycle) Lifecycle {
	return &lifecycle{internal: lc}
}

type lifecycle struct {
	internal fx.Lifecycle
}

func (lc *lifecycle) OnStart(f func(context.Context) error) {
	lc.internal.Append(fx.Hook{
		OnStart: f,
	})
}

func (lc *lifecycle) OnStop(f func(context.Context) error) {
	lc.internal.Append(fx.Hook{
		OnStop: f,
	})
}

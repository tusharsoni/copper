package chttp

import (
	"github.com/tusharsoni/copper"
)

type Router struct {
	copper.Module

	Routes []Route `group:"chttp/routes,flatten"`
}

func NewRouter(routes []Route) Router {
	return Router{Routes: routes}
}

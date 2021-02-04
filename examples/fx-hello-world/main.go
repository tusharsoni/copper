package main

import (
	"github.com/tusharsoni/copper"
	"github.com/tusharsoni/copper/cacl"
	"github.com/tusharsoni/copper/cauth"
	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
	"github.com/tusharsoni/copper/cqueue"
	"github.com/tusharsoni/copper/csql"
)

func main() {
	var app = copper.New()

	app.AddConfigs(
		csql.Config{
			Host: "127.0.0.1",
			Port: 5432,
			Name: "copper_ex",
			User: "postgres",
		},
	)

	app.AddModules(
		clogger.New,
		chttp.New,
		csql.New,
		cauth.New,
		cqueue.New,
		cacl.New,

		NewLogTask,
		NewAppRouter,
	)

	app.Start(
		cqueue.StartBackgroundWorkers,
		chttp.StartServer,
	)
}

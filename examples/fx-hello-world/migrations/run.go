package main

import (
	"github.com/tusharsoni/copper"
	"github.com/tusharsoni/copper/cacl"
	"github.com/tusharsoni/copper/cauth"
	"github.com/tusharsoni/copper/clogger"
	"github.com/tusharsoni/copper/cqueue"
	"github.com/tusharsoni/copper/csql"
)

func main() {
	var app = copper.New()

	app.AddModules(
		clogger.New,
		csql.New,
	)

	app.AddConfigs(
		csql.Config{
			Host: "127.0.0.1",
			Port: 5432,
			Name: "copper_ex",
			User: "postgres",
		},
	)

	app.Run(
		cacl.RunMigrations,
		cauth.RunMigrations,
		cqueue.RunMigrations,
	)
}

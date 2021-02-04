package main

import (
	"github.com/tusharsoni/copper"
	"github.com/tusharsoni/copper/cauth"
	cauthphone "github.com/tusharsoni/copper/cauth/phone"
	"github.com/tusharsoni/copper/clogger"
	"github.com/tusharsoni/copper/csql"
)

func main() {
	var app = copper.New()

	app.AddConfigs(csql.Config{
		Host:     "127.0.0.1",
		Port:     5432,
		Name:     "copper",
		User:     "postgres",
	})

	app.AddModules(
		clogger.New,
		csql.New,
	)

	app.Run(
		cauth.RunMigrations,
		cauthphone.RunMigrations,
	)
}

package main

import (
	"github.com/tusharsoni/copper"
	"github.com/tusharsoni/copper/cauth"
	cauthemail "github.com/tusharsoni/copper/cauth/email"
	cauthphone "github.com/tusharsoni/copper/cauth/phone"
	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
	"github.com/tusharsoni/copper/csql"
)

func main() {
	var app = copper.New()

	app.AddModules(
		clogger.New,
		//ctexter.FxLogger,
		csql.New,
		cauth.New,
		cauthphone.New,
		chttp.New,
	)

	app.AddConfigs(
		cauthemail.GetDefaultConfig(),
		csql.Config{
			Host: "127.0.0.1",
			Port: 5432,
			Name: "copper",
			User: "postgres",
		},
	)

	app.Start(
		chttp.StartServer,
	)
}

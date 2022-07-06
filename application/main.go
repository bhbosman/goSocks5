package main

import (
	"github.com/bhbosman/goFxApp"
	"github.com/bhbosman/goSocks5/providers"
)

func main() {
	app := goFxApp.NewFxMainApplicationServices(
		"Socks App",
		false,
		providers.ProvideSocksConnection(),
	)
	if app.FxApp.Err() != nil {
		println(app.FxApp.Err().Error())
		return
	}
	app.RunTerminalApp()
}

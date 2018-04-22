package main

import (
	"flag"
	"fmt"
	"sidekick/tmatrix/app"
	"sidekick/tmatrix/utils"
	"xframe/server"
)

var (
	config = flag.String("c", "", "configuration file path")
)

func main() {
	var app app.Config
	err := config.LoadConfigFromFileV2(&app, *config)

	//TODO  use errd
	if err != nil {
		panic(fmt.Sprintf("Load configuration error: %v", err))
	}

	//init http service conf
	utils.InitHttp(app.HttpConf)

	//init log
	utils.InitLog(app.LogConf)

	//start service
	if err = server.RunHTTP(utils.GetAddress(), utils.GetPort()); err != nil {
		panic(fmt.Sprintf("run tmatric service error: %v", err))
	}
}

package main

import (
	"flag"
	"fmt"
	"sidekick/tmatrix/app"
	"sidekick/tmatrix/logic/api"
	"sidekick/tmatrix/logic/conn"
	_ "sidekick/tmatrix/logic/service"
	"sidekick/tmatrix/utils"
	"xframe/config"
	"xframe/server"
)

var (
	conf = flag.String("c", "", "configuration file path")
)

func main() {
	var app app.Config
	err := config.LoadConfigFromFileV2(&app, *conf)

	//TODO  use errd
	if err != nil {
		panic(fmt.Sprintf("Load configuration error: %v", err))
	}

	//init http service conf
	utils.InitHttp(app.HttpConf)

	//init log
	utils.InitLog(app.LogConf)

	//init conn manager
	conn.Init(app)

	//init api manager
	api.Init(app)

	//start service
	if err = server.RunHTTP(utils.GetAddr(), utils.GetPort()); err != nil {
		panic(fmt.Sprintf("run tmatric service error: %v", err))
	}
}

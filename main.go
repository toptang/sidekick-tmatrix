package main

import (
	"flag"
	"fmt"
	"sidekick/tmatrix/app"
	"sidekick/tmatrix/logic/api"
	"sidekick/tmatrix/logic/conn"
	_ "sidekick/tmatrix/logic/service"
	"sidekick/tmatrix/utils"
	"xframe/cmd"
	"xframe/config"
	"xframe/server"
)

var (
	conf = flag.String("c", "", "configuration file path")
)

func main() {
	//init commandLine
	cmd.ParseCommand()
	cmd.DumpCommand()

	//init configuration
	var appConf app.Config
	err := config.LoadConfigFromFileV2(&appConf, *conf)

	//TODO  use errd
	if err != nil {
		panic(fmt.Sprintf("Load configuration error: %v", err))
	}

	//init config
	app.InitApp(appConf)

	//init conn manager
	conn.Init(appConf)

	//init api manager
	api.Init(appConf)

	//start service
	if err = server.RunHTTP(utils.GetHttpAddr(), utils.GetHttpPort()); err != nil {
		panic(fmt.Sprintf("run tmatric service error: %v", err))
	}
}

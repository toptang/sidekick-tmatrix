package app

import "sidekick/tmatrix/utils"

type Config struct {
	HttpConf     utils.HttpConfig `json:"http"`
	LogConf      utils.LogConfig  `json:"log"`
	UpstreamConf struct {
		OkexConf utils.OkexConfig `json:"okex"`
	} `json:"upstream"`
}

func InitApp(app Config) {
	//init http service conf
	utils.InitHttp(app.HttpConf)

	//init log
	utils.InitLog(app.LogConf)
}

package api

import (
	"sidekick/tmatrix/app"
	"sidekick/tmatrix/logic/api/okexapi"
	"sidekick/tmatrix/utils"
	"sync"
	"xframe/log"
)

var (
	GApi sync.Map
)

type BaseApi interface {
	Start(string, string, string, int, string) error
}

func Init(config app.Config) {
	if config.UpstreamConf.OkexConf.Enabled {
		utils.InitOkexConfig(config.UpstreamConf.OkexConf)
		log.DEBUG("load okex api")
		okexApi := okexapi.NewOkexApi()
		GApi.Store("okex", okexApi)
	}
}

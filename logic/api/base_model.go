package api

import (
	"sidekick/tmatrix/app"
	okexapi "sidekick/tmatrix/logic/api/okex"
	"sync"
	"xframe/log"
)

var (
	GApi sync.Map
)

type BaseApi interface {
	Start(string, string)
}

func Init(config app.Config) {
	if config.UpstreamConf.OkexConf.Enabled {
		log.DEBUG("load okex api")
		okexApi := okexapi.NewOkexApi()
		GApi.Store("okex", okexApi)
	}
}

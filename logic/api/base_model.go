package api

import (
	"sidekick/tmatrix/app"
	"sidekick/tmatrix/logic/api/okexapi"
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
		log.DEBUG("load okex api")
		okexApi := okexapi.NewOkexApi()
		GApi.Store("okex", okexApi)
	}
}

package api

import (
	"sidekick/tmatrix/app"
	okexapi "sidekick/tmatrix/logic/api/okex"
	"sync"
)

var (
	GApi sync.Map
)

type BaseApi interface {
	Start(string, string)
}

func Init(config app.Config) {
	if config.UpstreamConf.OkexConf.Enabled {
		okexApi := okexapi.NewOkexApi()
		GApi.Store("okex", okexApi)
	}
}

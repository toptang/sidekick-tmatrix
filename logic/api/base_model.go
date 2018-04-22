package api

import (
	"sidekick/tmatrix/app"
	"sync"
)

var (
	GApi sync.Map
)

type BaseApi interface {
	Start(string, string)
}

func Init(config app.Config) {
	if config.ApiConf.OkexConf.Enabled {
		okexApi := NewOkexApi()
		GApi.Store("okex", okexApi)
	}
}

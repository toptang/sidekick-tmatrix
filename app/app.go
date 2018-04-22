package app

import "sidekick/tmatrix/utils"

type Config struct {
	HttpConf     utils.HttpConfig `json:"http"`
	LogConf      utils.LogConfig  `json:"log"`
	UpstreamConf struct {
		OkexConf OkexConfig `json:"okex"`
	} `json:"upstream"`
}

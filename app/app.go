package app

import "sidekick/tmatrix/utils"

type Config struct {
	HttpConf     utils.HttpConfig `json:"http"`
	LogConf      utils.LogConfig  `json:"log"`
	UpstreamConf struct {
		OkexConf utils.OkexConfig `json:"okex"`
	} `json:"upstream"`
}

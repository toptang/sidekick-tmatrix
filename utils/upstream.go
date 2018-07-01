package utils

type OkexConfig struct {
	Enabled   bool   `json:"enabled"` //used for register okex upstream websocket and connection manager
	ApiKey    string `json:"api_key"`
	ApiSecret string `json:"api_secret"`
}

var (
	okexConfig OkexConfig
)

func InitOkexConfig(OkexConf OkexConfig) {
	okexConfig = OkexConf
	if okexConfig.ApiKey == EMPTY_STR ||
		okexConfig.ApiSecret == EMPTY_STR {
		panic("okex config error")
	}
}

func GetOkexKey() string {
	return okexConfig.ApiKey
}

func GetOkexSecret() string {
	return okexConfig.ApiSecret
}

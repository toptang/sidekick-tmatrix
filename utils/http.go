package utils

type HttpConfig struct {
	Addr string `json:"address"`
	Port int    `json:"port"`
}

var (
	httpConfig *HttpConfig
)

func InitHttp(httpConf HttpConfig) {
	httpConfig = &httpConf
	if httpConfig.Addr == "" ||
		httpConfig.Port == 0 {
		panic("http service config error")
	}
}

func GetAddr() string {
	return httpConfig.Addr
}

func GetPort() int {
	return httpConfig.Port
}

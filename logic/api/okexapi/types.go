package okexapi

import "time"

//basic fixed config
var (
	//websocket
	OKEX_SERVER = "wss://real.okex.com:10440"
	OKEX_URI    = "/websocket/okexapi"
	PROTO_WS    = "ws"
	PROTO_WSS   = "wss"

	//REST api
	OKEX_API_SERVER  = "https://www.okex.com/api/v1"
	OKEX_TRADE_URI   = "/future_trade"
	OKEX_UNTRADE_URI = "/future_cancel"
)

var (
	DEFAULT_PERIOD = "this_week"
	DEFAULT_DEPTH  = 5

	HEALTH_CHECK_TIME        = 30 * time.Second
	HEALTH_CHECK_RETRY       = 3
	HEALTH_GAP         int64 = 40
)

//protocol uri
var (
	OKEX_OB     = "ok_sub_futureusd_%s_depth_%s_%d"
	OKEX_TRADE  = "ok_sub_futureusd_%s_trade_%s"
	OKEX_TICKER = "ok_sub_futureusd_%s_ticker_%s"
)

//-----------------
type PingPongRequest struct {
	Event string `json:"event"`
}

//-----------------
type Param struct {
	ApiKey string `json:"api_key"`
	Sign   string `json:"sign"`
}

type DataRequest struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	Params  Param  `json:"params"`
}

//-----------------
type PingPongResponse struct {
	Event string `json:"event"`
}

//-----------------
type DataResponse struct {
	Channel string `json:"channel"`
	Data    Quote  `json:"data"`
}

type Quote struct {
	Ts   int64       `json:"timestamp"`
	Asks interface{} `json:"asks"`
	Bids interface{} `json:"bids"`
}

//-----------------
type QuotePush struct {
	Ts       int64       `json:"ts"`
	Asks     interface{} `json:"asks"`
	Bids     interface{} `json:"bids"`
	Contract string      `json:"symbol"`
	Type     string      `json:"type"`
}

type Option struct {
	Period string `json:"period"`
	Depth  int    `json:"depth"`
}

type DataPush struct {
	Msg      string      `json:"msg"`
	Market   string      `json:"market"`
	Table    string      `json:"table"`
	Contract string      `json:"symbol"`
	Type     string      `json:"type"`
	Data     []QuotePush `json:"data"`
	Optional Option      `json:"optional"`
	Channel  string      `json:"channel"`
}

//-----------------
type ChannelRes struct {
	Channel string `json:"channel"`
}

//-----------------
type DataCommonRes struct {
	Channel string      `json:"channel"`
	Data    interface{} `json:"data"`
}

type DataCommonPush struct {
	Msg      string        `json:"msg"`
	Market   string        `json:"market"`
	Table    string        `json:"table"`
	Contract string        `json:"symbol"`
	Type     string        `json:"type"`
	Data     []interface{} `json:"data"`
	Optional Option        `json:"optional"`
	Channel  string        `json:"channel"`
}

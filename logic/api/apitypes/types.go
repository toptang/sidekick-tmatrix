package apitypes

import "time"

var (
	API_DATA_ORDERBOOK = "orderbook"
	API_DATA_LOGIN     = "login"
	API_DATA_TICKER    = "ticker"
)

//health check
var (
	HEALTH_CHECK_TIME        = 30 * time.Second
	HEALTH_CHECK_RETRY       = 3
	HEALTH_GAP         int64 = 40
)

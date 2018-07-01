package okexapi

import (
	"encoding/json"
	"fmt"
	"sidekick/tmatrix/utils"
	"xframe/log"
	"xframe/server/websocket"
)

//-------------------------
//push to client
func (this *OkexApi) generateQuotePushData(dataRes []DataResponse, depth int) ([]byte, error) {
	return this.generateCommonQuoteData("quote", dataRes, depth)
}

func (this *OkexApi) generateLoginPushData(dataCommonRes []DataCommonRes, pushTable string) ([]byte, error) {
	return this.generateCommonData("trade", dataCommonRes, pushTable)
}

func (this *OkexApi) generateTickerPushData(dataCommonRes []DataCommonRes) ([]byte, error) {
	return this.generateCommonData("quote", dataCommonRes, this.Table)
}

func (this *OkexApi) generateSpotTickerPushData(dataCommonRes []DataCommonRes) ([]byte, error) {
	return this.generateCommonData("spot_quote", dataCommonRes, this.Table)
}

func (this *OkexApi) generateSpotOrderbookPushData(dataRes []DataResponse, depth int) ([]byte, error) {
	return this.generateCommonQuoteData("spot_quote", dataRes, depth)
}

//---------------------------------------------
// request to okex
// 合约交易
func (this *OkexApi) SendPingPong(wsConn *websocket.Conn) {
	ppReq := PingPongRequest{
		Event: "ping",
	}
	buf, _ := json.Marshal(ppReq)
	wsConn.Write(buf)
}

func (this *OkexApi) SendFutureUsd(wsConn *websocket.Conn, depth int) {
	//signature
	params := map[string]string{
		"api_key": utils.GetOkexKey(),
	}
	sign := this.Sign(params, utils.GetOkexSecret())
	param := Param{
		ApiKey: utils.GetOkexKey(),
		Sign:   sign,
	}
	dataReq := DataRequest{
		Event:   "addChannel",
		Channel: fmt.Sprintf(OKEX_OB, this.Symbol, this.Type, depth),
		Params:  param,
	}
	buf, _ := json.Marshal(dataReq)
	wsConn.Write(buf)
}

func (this *OkexApi) SendFutureUsdTicker(wsConn *websocket.Conn) {
	params := map[string]string{
		"api_key": utils.GetOkexKey(),
	}
	sign := this.Sign(params, utils.GetOkexSecret())
	param := Param{
		ApiKey: utils.GetOkexKey(),
		Sign:   sign,
	}
	dataReq := DataRequest{
		Event:   "addChannel",
		Channel: fmt.Sprintf(OKEX_TICKER, this.Symbol, this.Type),
		Params:  param,
	}
	log.DEBUG(dataReq)
	buf, _ := json.Marshal(dataReq)
	wsConn.Write(buf)
}

func (this *OkexApi) SendLogin(wsConn *websocket.Conn) {
	//signature
	params := map[string]string{
		"api_key": utils.GetOkexKey(),
	}
	sign := this.Sign(params, utils.GetOkexSecret())
	param := Param{
		ApiKey: utils.GetOkexKey(),
		Sign:   sign,
	}
	dataReq := DataRequest{
		Event:  "login",
		Params: param,
	}
	buf, _ := json.Marshal(dataReq)
	wsConn.Write(buf)
}

//---------------------------------------------
// request to okex
// 现货交易
func (this *OkexApi) SendSpotTicker(wsConn *websocket.Conn) {
	params := map[string]string{
		"api_key": utils.GetOkexKey(),
	}
	sign := this.Sign(params, utils.GetOkexSecret())
	param := Param{
		ApiKey: utils.GetOkexKey(),
		Sign:   sign,
	}
	dataReq := DataRequest{
		Event:   "addChannel",
		Channel: fmt.Sprintf(OKEX_SPOT_TICKER, this.Symbol),
		Params:  param,
	}
	buf, _ := json.Marshal(dataReq)
	wsConn.Write(buf)
}

func (this *OkexApi) SendSpotOrderbook(wsConn *websocket.Conn, depth int) {
	params := map[string]string{
		"api_key": utils.GetOkexKey(),
	}
	sign := this.Sign(params, utils.GetOkexSecret())
	param := Param{
		ApiKey: utils.GetOkexKey(),
		Sign:   sign,
	}
	dataReq := DataRequest{
		Event:   "addChannel",
		Channel: fmt.Sprintf(OKEX_SPOT_OB, this.Symbol, depth),
		Params:  param,
	}
	buf, _ := json.Marshal(dataReq)
	wsConn.Write(buf)
}

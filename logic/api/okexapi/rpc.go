package okexapi

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"sidekick/tmatrix/utils"
	"sort"
	"strings"
	"xframe/log"
	"xframe/server/websocket"
)

//-------------------------
//push
func (this *OkexApi) generateQuotePushData(dataRes []DataResponse, market string, contract string, table string, ttype string, depth int) ([]byte, error) {
	var (
		dataPush DataPush
	)
	dataPush.Msg = "quote"
	dataPush.Market = market
	dataPush.Table = table
	dataPush.Contract = contract
	dataPush.Type = ttype
	dataPush.Optional = Option{
		Period: "",
		Depth:  depth,
	}
	dataPush.Data = make([]QuotePush, 0)
	for _, data := range dataRes {
		var tmpQuotePush QuotePush
		tmpQuotePush.Ts = data.Data.Ts
		tmpQuotePush.Asks = data.Data.Asks
		tmpQuotePush.Bids = data.Data.Bids
		//tmpQuotePush.Contract = contract
		//tmpQuotePush.Type = ttype
		dataPush.Data = append(dataPush.Data, tmpQuotePush)
		//dataPush.Channel = data.Channel
	}
	buf, err := json.Marshal(dataPush)
	return buf, err
}

func (this *OkexApi) generateLoginPushData(dataCommonRes []DataCommonRes, market string, contract string, table string, ttype string) ([]byte, error) {
	var (
		dataCommonPush DataCommonPush
	)
	dataCommonPush.Msg = "trade"
	dataCommonPush.Table = table
	dataCommonPush.Contract = contract
	dataCommonPush.Type = ttype
	dataCommonPush.Data = make([]interface{}, 0)
	//structure抽象
	for _, data := range dataCommonRes {
		dataCommonPush.Data = append(dataCommonPush.Data, data.Data)
	}
	buf, err := json.Marshal(dataCommonPush)
	return buf, err
}

func (this *OkexApi) generateTickerPushData(dataCommonRes []DataCommonRes, market string, contract string, table string, ttype string) ([]byte, error) {
	var (
		dataCommonPush DataCommonPush
	)
	dataCommonPush.Msg = "quote"
	dataCommonPush.Market = market
	dataCommonPush.Table = table
	dataCommonPush.Contract = contract
	dataCommonPush.Type = ttype
	dataCommonPush.Data = make([]interface{}, 0)
	//TODO structure抽象
	for _, data := range dataCommonRes {
		dataCommonPush.Data = append(dataCommonPush.Data, data.Data)
	}
	buf, err := json.Marshal(dataCommonPush)
	return buf, err
}

//---------------------------------------------
// request
func (this *OkexApi) SendPingPong(wsConn *websocket.Conn) {
	ppReq := PingPongRequest{
		Event: "ping",
	}
	buf, _ := json.Marshal(ppReq)
	wsConn.Write(buf)
}

func (this *OkexApi) SendFutureUsd(contract string, ttype string, wsConn *websocket.Conn, depth int) {
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
		Channel: fmt.Sprintf(OKEX_OB, contract, ttype, depth),
		Params:  param,
	}
	buf, _ := json.Marshal(dataReq)
	wsConn.Write(buf)
}

func (this *OkexApi) SendFutureUsdTicker(contract string, ttype string, wsConn *websocket.Conn) {
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
		Channel: fmt.Sprintf(OKEX_TICKER, contract, ttype),
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

func (this *OkexApi) Sign(params map[string]string, api_secret string) string {
	var (
		keyLst     = make([]string, 0)
		sortParams string
	)
	for key, _ := range params {
		keyLst = append(keyLst, key)
	}
	sort.Strings(keyLst)
	for _, key := range keyLst {
		sortParams += key + "=" + params[key] + "&"
	}
	sortParams += "secret_key=" + api_secret
	h := md5.New()
	io.WriteString(h, sortParams)
	sign := strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil)))
	return sign
}

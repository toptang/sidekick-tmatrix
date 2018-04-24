package okexapi

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"sidekick/tmatrix/utils"
	"sort"
	"strings"
	"xframe/server/websocket"
)

func (this *OkexApi) SendPingPong(wsConn *websocket.Conn) {
	ppReq := PingPongRequest{
		Event: "ping",
	}
	buf, _ := json.Marshal(ppReq)
	wsConn.Write(buf)
}

func (this *OkexApi) SendFutureUsed(contract string, wsConn *websocket.Conn) {
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
		Channel: fmt.Sprintf(OKEX_OB, contract, DEFAULT_PERIOD, DEFAULT_DEPTH),
		Params:  param,
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

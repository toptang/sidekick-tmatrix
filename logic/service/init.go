package service

import (
	"encoding/json"
	"sidekick/tmatrix/logic/service/handlers"
	"sidekick/tmatrix/logic/service/requests"
	"sidekick/tmatrix/logic/service/response"
	"sidekick/tmatrix/logic/service/svcerr"
	"time"
	"xframe/handler/websocket_handler"
	"xframe/log"
	"xframe/server"
	"xframe/server/websocket"
)

var (
	wsHandlers = map[string]struct {
		Handler func(chan []byte, interface{}, interface{})
		Timeout time.Duration
	}{
		"sub":   {Handler: handlers.SubMsg, Timeout: 0},
		"unsub": {Handler: handlers.UnsubMsg, Timeout: 0},
	}
)

func init() {
	for name, wsHandler := range wsHandlers {
		websocket_handler.RegisterWsTaskHandle(name, websocket_handler.WsTaskFunc(wsHandler.Handler), wsHandler.Timeout)
	}
	server.RouteWs = RouteWs
}

func _wsOnMessages(ws *websocket.Conn) {
	for {
		var (
			buf       []byte             //binary stream
			pubSubReq requests.PubSubReq //proto
		)
		err := websocket.Message.Receive(ws, &buf)
		if err != nil {
			log.WARNF("[_wsOnMessages]websocket client receive error: %v", err)
			return
		}
		//Decode proto
		json.Unmarshal(buf, &pubSubReq)
		if pubSubReq.CheckRouter() {
			log.ERRORF("[_wsOnMessages]router action not found: %s", pubSubReq.Msg)
			response.DoBaseResponse(pubSubReq.Msg, svcerr.ROUTE_ERROR, pubSubReq.Uuid, ws)
			continue
		}
		task, err := websocket_handler.NewWsTask(pubSubReq.Msg)
		if err != nil {
			log.ERRORF("[_wsOnMessages]not found task for %s", pubSubReq.Msg)
			response.DoBaseResponse(pubSubReq.Msg, svcerr.INTERNAL_ERROR, pubSubReq.Uuid, ws)
			continue
		}
		//FIXME: no tracing
		res, err := task.Run(pubSubReq, ws)
		if err != nil {
			log.ERRORF("[_wsOnMessages]%s task run fail", pubSubReq.Msg)
			response.DoBaseResponse(pubSubReq.Msg, svcerr.INTERNAL_ERROR, pubSubReq.Uuid, ws)
			continue
		}
		if res != nil {
			go response.DoDataResponse(res, ws)
		}
	}
}

func RouteWs(ws *websocket.Conn) {
	_wsOnMessages(ws)
}

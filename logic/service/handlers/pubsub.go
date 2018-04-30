package handlers

import (
	"sidekick/tmatrix/logic/api"
	okexapi "sidekick/tmatrix/logic/api/okex"
	"sidekick/tmatrix/logic/conn"
	"sidekick/tmatrix/logic/service/requests"
	"sidekick/tmatrix/logic/service/response"
	"sidekick/tmatrix/logic/service/svcerr"
	"xframe/log"
	"xframe/server/websocket"
)

func SubMsg(msgChan chan []byte, msg interface{}, ws interface{}) {
	//define msg to pubSubReq
	req := msg.(requests.PubSubReq)
	//check params
	if req.CheckParams() {
		log.ERRORF("[sub_msg]params error")
		response.DoBaseResponse(req.Msg, svcerr.INPUT_ERROR, req.Uuid, ws.(*websocket.Conn))
		return
	}
	//get manager by market
	if manager, ok := conn.GConn.Load(req.Data.Market); !ok {
		log.ERRORF("[sub_msg]not found %s manager", req.Data.Market)
		response.DoBaseResponse(req.Msg, svcerr.CONN_NOT_FOUND_ERROR, req.Uuid, ws.(*websocket.Conn))
		return
	} else {
		//register client conn
		switch req.Data.Market {
		case "okex":
			//FIXME sync.map interface{} != defined customized interface{}
			connManager := manager.(*conn.OKEXManager)
			connManager.RegisterConn(ws.(*websocket.Conn), req.Data.Symbol, req.Data.Table)
			cliLst := connManager.DumpConns(req.Data.Symbol, req.Data.Table)
			log.DEBUG("[sub_msg]cliLst", cliLst)
		}
	}
	if apicli, ok := api.GApi.Load(req.Data.Market); !ok {
		log.ERRORF("[sub_msg]not found %s api", req.Data.Market)
		response.DoBaseResponse(req.Msg, svcerr.API_NOT_FOUND_ERROR, req.Uuid, ws.(*websocket.Conn))
		return
	} else {
		//start upstream client
		switch req.Data.Market {
		case "okex":
			baseApi := apicli.(*okexapi.OkexApi)
			go baseApi.Start(req.Data.Symbol, req.Data.Table)
		}
	}
	log.INFOF("[sub_msg]complete sub in market: %s, symbol: %s, table: %s", req.Data.Market, req.Data.Symbol, req.Data.Table)
	response.DoBaseResponse(req.Msg, svcerr.SUCCESS, req.Uuid, ws.(*websocket.Conn))
	msgChan <- nil
	return
}

func UnsubMsg(msgChan chan []byte, msg interface{}, ws interface{}) {
	//define msg to pubSubReq
	req := msg.(requests.PubSubReq)
	//check params
	if req.CheckParams() {
		log.ERRORF("[unsub_msg]params error")
		response.DoBaseResponse(req.Msg, svcerr.INPUT_ERROR, req.Uuid, ws.(*websocket.Conn))
		return
	}
	//get manager by market
	if manager, ok := conn.GConn.Load(req.Data.Market); !ok {
		log.ERRORF("[unsub_msg]not found %s manager", req.Data.Market)
		response.DoBaseResponse(req.Msg, svcerr.CONN_NOT_FOUND_ERROR, req.Uuid, ws.(*websocket.Conn))
		return
	} else {
		//unregister client conn
		//connManager := manager.(conn.ConnManager)
		switch req.Data.Market {
		case "okex":
			connManager := manager.(*conn.OKEXManager)
			go connManager.UnRegisterConn(ws.(*websocket.Conn), req.Data.Symbol, req.Data.Table)
		}
	}
	log.INFOF("[unsub_msg]complete unsub in market: %s, symbol: %s, table: %s", req.Data.Market, req.Data.Symbol, req.Data.Table)
	response.DoBaseResponse(req.Msg, svcerr.SUCCESS, req.Uuid, ws.(*websocket.Conn))
	msgChan <- nil
	return
}

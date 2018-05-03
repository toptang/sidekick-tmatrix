package okexconn

import (
	"fmt"
	"sidekick/tmatrix/logic/conn/conntypes"
	"xframe/log"
	"xframe/server/websocket"
)

type OpReq struct {
	client  *conntypes.UpstreamClient
	msg     string
	retChan chan map[string]*conntypes.UpstreamClient
}

/*
 * OKEX连接管理
 * - 添加客户端连接(SUB)
 * - 删除客户端连接(UNSUB)
 * - 收上游推送消息 -> 广播客户端
 */
type OKEXManager struct {
	inputChan chan []byte
	ClientLst map[string]map[string]*conntypes.UpstreamClient //key: contract|table
	opChan    chan OpReq
}

func NewOKEXManager() *OKEXManager {
	om := new(OKEXManager)
	om.ClientLst = make(map[string]map[string]*conntypes.UpstreamClient) //sub-key:remoteaddr
	om.opChan = make(chan OpReq)
	go om.RunOp()
	return om
}

func (this *OKEXManager) register(key string, client *conntypes.UpstreamClient) {
	remoteAddr := client.RemoteAddr
	if okexCliMap, ok := this.ClientLst[key]; !ok {
		this.ClientLst[key] = make(map[string]*conntypes.UpstreamClient)
		this.ClientLst[key][remoteAddr] = client
	} else {
		//allow to duplicately sub
		okexCliMap[remoteAddr] = client
	}
	log.DEBUGF("[okex_manager]after register: %v", this.ClientLst)
}

func (this *OKEXManager) unregister(key string, client *conntypes.UpstreamClient) {
	remoteAddr := client.RemoteAddr
	if okexCliMap, ok := this.ClientLst[key]; ok {
		delete(okexCliMap, remoteAddr)
	}
	log.DEBUGF("[okex_manager]after unregister: %v", this.ClientLst)

}

func (this *OKEXManager) RunOp() {
	log.DEBUG("[okex_manager]start run op")
	for {
		select {
		case req := <-this.opChan:
			var (
				contract string = req.client.Contract
				table    string = req.client.Table
				ttype    string = req.client.Type
				key      string = this.getKey(contract, table, ttype)
			)

			switch req.msg {
			case conntypes.REGISTER_ROUTE:
				log.DEBUG("[okex_manager]register conn")
				this.register(key, req.client)
			case conntypes.UNREGISTER_ROUTE:
				log.DEBUG("[okex_manager]unregister conn")
				this.unregister(key, req.client)
			case conntypes.DUMP_ROUTE:
				req.retChan <- this.ClientLst[key]
			}
		}
	}
}

func (this *OKEXManager) RegisterConn(ws *websocket.Conn, contract string, table string, ttype string) {
	log.DEBUGF("[okex_manager]register conn for contract %s, table %s and type %s", contract, table, ttype)
	okexClient := conntypes.NewUpstreamClient(ws.RemoteAddr().String(), contract, table, ttype, ws)
	var opReq = OpReq{
		client: okexClient,
		msg:    conntypes.REGISTER_ROUTE,
	}
	go func() {
		this.opChan <- opReq
	}()
}

func (this *OKEXManager) UnRegisterConn(ws *websocket.Conn, contract string, table string, ttype string) {
	log.DEBUGF("[okex_manager]unregister conn for contract %s, table %s and type %s", contract, table, ttype)
	okexClient := conntypes.NewUpstreamClient(ws.RemoteAddr().String(), contract, table, ttype, ws)
	var opReq = OpReq{
		client: okexClient,
		msg:    conntypes.UNREGISTER_ROUTE,
	}
	go func() {
		this.opChan <- opReq
	}()
}

func (this *OKEXManager) DumpConns(contract string, table string, ttype string) map[string]*conntypes.UpstreamClient {
	okexClient := conntypes.NewUpstreamClient("", contract, table, ttype, nil)
	rCh := make(chan map[string]*conntypes.UpstreamClient)
	var opReq = OpReq{
		client:  okexClient,
		msg:     conntypes.DUMP_ROUTE,
		retChan: rCh,
	}
	go func(data OpReq) {
		this.opChan <- data
	}(opReq)
	select {
	case res := <-rCh:
		log.DEBUGF("[dump_conns]dump client list for contract %s, table %s, type %s: %v", contract, table, ttype, res)
		return res
	}
}

func (this *OKEXManager) getKey(contract string, table string, ttype string) string {
	return fmt.Sprintf("%s|%s|%s", contract, table, ttype)
}

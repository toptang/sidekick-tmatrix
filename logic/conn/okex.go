package conn

import (
	"fmt"
	"xframe/server/websocket"
)

type OpReq struct {
	client  *OKEXClient
	msg     string
	retChan chan map[string]*OKEXClient
}

/*
 * OKEX连接管理
 * - 添加客户端连接(SUB)
 * - 删除客户端连接(UNSUB)
 * - 收上游推送消息 -> 广播客户端
 */
type OKEXManager struct {
	inputChan chan []byte

	ClientLst map[string]map[string]*OKEXClient //key: contract|table

	opChan chan OpReq
}

//TODO: try to use sync.Map
//      abstract table struct
func NewOKEXManager() *OKEXManager {
	om := new(OKEXManager)
	om.ClientLst = make(map[string]map[string]*OKEXClient) //sub-key:remoteaddr
	om.opChan = make(chan OpReq)
	om.RunOp()
	return om
}

func (this *OKEXManager) RunOp() {
	for {
		select {
		case req := <-this.opChan:
			contract := req.client.Contract
			table := req.client.Table
			remoteAddr := req.client.RemoteAddr
			key := this.getKey(contract, table)

			switch req.msg {
			case REGISTER_ROUTE:
				if okexCliMap, ok := this.ClientLst[key]; !ok {
					this.ClientLst[key] = make(map[string]*OKEXClient)
					this.ClientLst[key][remoteAddr] = req.client
				} else {
					if _, ok := okexCliMap[remoteAddr]; !ok {
						this.ClientLst[key][remoteAddr] = req.client
					}
				}
			case UNREGISTER_ROUTE:
				if okexCliMap, ok := this.ClientLst[key]; ok {
					delete(okexCliMap, remoteAddr)
				}
			case DUMP_ROUTE:
				req.retChan <- this.ClientLst[key]
			}
		}
	}
}

func (this *OKEXManager) RegisterConn(ws *websocket.Conn, contract string, table string) {
	okexClient := NewOKEXClient(ws.RemoteAddr().String(), contract, table, ws)
	var opReq = OpReq{
		client: okexClient,
		msg:    REGISTER_ROUTE,
	}
	go func() {
		this.opChan <- opReq
	}()
}

func (this *OKEXManager) UnRegisterConn(ws *websocket.Conn, contract string, table string) {
	okexClient := NewOKEXClient(ws.RemoteAddr().String(), contract, table, ws)
	var opReq = OpReq{
		client: okexClient,
		msg:    UNREGISTER_ROUTE,
	}
	go func() {
		this.opChan <- opReq
	}()
}

func (this *OKEXManager) DumpConns(contract string, table string) map[string]*OKEXClient {
	okexClient := NewOKEXClient("", contract, table, nil)
	rCh := make(chan map[string]*OKEXClient)
	var opReq = OpReq{
		client:  okexClient,
		msg:     DUMP_ROUTE,
		retChan: rCh,
	}
	go func() {
		this.opChan <- opReq
	}()
	select {
	case res := <-rCh:
		return res
	}
}

func (this *OKEXManager) getKey(contract string, table string) string {
	return fmt.Sprintf("%s|%s", contract, table)
}

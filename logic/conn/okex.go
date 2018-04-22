package conn

import (
	"fmt"
	"xframe/server/websocket"
)

type OpReq struct {
	client  *OKEKClient
	msg     string
	retChan chan interface{}
}

/*
 * OKEK连接管理
 * - 添加客户端连接(SUB)
 * - 删除客户端连接(UNSUB)
 * - 收上游推送消息 -> 广播客户端
 */
type OKEKManager struct {
	inputChan chan []byte

	ClientLst map[string]map[string]*OKEKClient //key: contract|table

	opChan chan OpReq
}

//TODO: try to use sync.Map
//      abstract table struct
func NewOKEKManager() *OKEKManager {
	om := new(OKEKManager)
	om.ClientLst = make(map[string]map[string]*OKEKClient) //sub-key:remoteaddr
	om.opChan = make(chan OpReq)
	om.RunOp()
	return om
}

func (this *OKEKManager) RunOp() {
	for {
		select {
		case req := <-this.opChan:
			contract := req.client.Contract
			table := req.client.Table
			remoteAddr = req.client.RemoteAddr
			key := this.getKey(contract, table)

			switch req.msg {
			case REGISTER_ROUTE:
				if okekCliMap, ok := this.ClientLst[key]; !ok {
					this.ClientLst[key] = make(map[string]*OKEKClient)
					this.ClientLst[key][remoteAddr] = req.client
				} else {
					if _, ok := okekCliMap[remoteAddr]; !ok {
						this.ClientLst[key][remoteAddr] = req.client
					}
				}
			case UNREGISTER_ROUTE:
				if okekCliMap, ok := this.ClientLst[key]; ok {
					delete(okekCliMap, remoteAddr)
				}
			case DUMP_ROUTE:
				req.retChan <- this.ClientLst[key]
			}
		}
	}
}

func (this *OKEKManager) RegisterConn(ws *websocket.Conn, contract string, table string) {
	okekClient := NewOKEKClient(ws.RemoteAddr(), contract, table, ws)
	var opReq = OpReq{
		client: okekClient,
		msg:    REGISTER_ROUTE,
	}
	go func() {
		this.opChan <- opReq
	}()
}

func (this *OKEKManager) UnRegisterConn(ws *websocket.Conn, contract string, table string) {
	okekClient := NewOKEKClient(ws.RemoteAddr, contract, table, ws)
	var opReq = OpReq{
		client: okekClient,
		msg:    UNREGISTER_ROUTE,
	}
	go func() {
		this.opChan <- opReq
	}()
}

func (this *OKEKManager) DumpConns(contract string, table string) map[string]*OKEKClient {
	okekClient := NewOKEKClient("", contract, table, nil)
	rCh := make(chan interface{})
	var opReq = OpReq{
		client:  okekClient,
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

func (this *OKEKManager) getKey(contract string, table string) string {
	return fmt.Sprintf("%s|%s", contract, table)
}

package conntypes

import "xframe/server/websocket"

var (
	REGISTER_ROUTE   = "register"
	UNREGISTER_ROUTE = "unregister"
	DUMP_ROUTE       = "dump"
)

//RemoteAddr + Contract
type UpstreamClient struct {
	RemoteAddr string
	Contract   string
	Table      string
	Type       string
	Conn       *websocket.Conn
}

func NewUpstreamClient(addr string, contract string, table string, ttype string, ws *websocket.Conn) *UpstreamClient {
	return &UpstreamClient{
		RemoteAddr: addr,
		Contract:   contract,
		Conn:       ws,
		Type:       ttype,
		Table:      table,
	}
}

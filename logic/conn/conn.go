package conn

import (
	"sidekick/tmatrix/app"
	"sync"
	"xframe/server/websocket"
)

var (
	GConn sync.Map
)

var (
	REGISTER_ROUTE   = "register"
	UNREGISTER_ROUTE = "unregister"
	DUMP_ROUTE       = "dump"
)

/*
 * ws conn manager
 */
type ConnManager interface {
	RegisterConn(*websocket.Conn, string, string)
	UnRegisterConn(*websocket.Conn, string, string)
	DumpConns(string) map[string]*OKEKClient
}

//RemoteAddr + Contract
type OKEKClient struct {
	RemoteAddr string
	Contract   string
	Table      string
	Conn       *websocket.Conn
}

func NewOKEKClient(addr string, contract string, table string, ws *websocket.Conn) *OKEKClient {
	return &OKEKClient{
		RemoteAddr: addr,
		Contract:   contract,
		Conn:       ws,
		Table:      table,
	}
}

//register all upstream manager
func Init(config app.Config) {
	if config.UpstreamConf.OkekConf.Enabled {
		okekManager := NewOKEKManager()
		GConn.Store("okek", okekManager)
	}
}

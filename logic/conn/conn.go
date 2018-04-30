package conn

import (
	"sidekick/tmatrix/app"
	"sync"
	"xframe/log"
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
	DumpConns(string, string) map[string]*OKEXClient
}

//RemoteAddr + Contract
type OKEXClient struct {
	RemoteAddr string
	Contract   string
	Table      string
	Conn       *websocket.Conn
}

func NewOKEXClient(addr string, contract string, table string, ws *websocket.Conn) *OKEXClient {
	return &OKEXClient{
		RemoteAddr: addr,
		Contract:   contract,
		Conn:       ws,
		Table:      table,
	}
}

//register all upstream manager
func Init(config app.Config) {
	if config.UpstreamConf.OkexConf.Enabled {
		log.DEBUG("load okex manager")
		okexManager := NewOKEXManager()
		GConn.Store("okex", okexManager)
	}
}

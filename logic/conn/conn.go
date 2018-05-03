package conn

import (
	"sidekick/tmatrix/app"
	"sidekick/tmatrix/logic/conn/conntypes"
	"sidekick/tmatrix/logic/conn/okexconn"
	"sync"
	"xframe/log"
	"xframe/server/websocket"
)

var (
	GConn sync.Map
)

/*
 * ws conn manager
 */
type ConnManager interface {
	RegisterConn(*websocket.Conn, string, string, string)
	UnRegisterConn(*websocket.Conn, string, string, string)
	DumpConns(string, string, string) map[string]*conntypes.UpstreamClient
	RunOp()
}

//register all upstream manager
func Init(config app.Config) {
	if config.UpstreamConf.OkexConf.Enabled {
		log.DEBUG("load okex manager")
		okexManager := okexconn.NewOKEXManager()
		GConn.Store("okex", okexManager)
	}
}

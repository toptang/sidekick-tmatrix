package okexapi

import (
	"encoding/json"
	"fmt"
	"sidekick/tmatrix/logic/conn"
	"sync"
	"time"
	"xframe/log"
	"xframe/server/websocket"
)

type OkexConn struct {
	WsConn     *websocket.Conn
	Health     bool
	HealthTime int64
}

func NewOkexConn(ws *websocket.Conn) *OkexConn {
	okc := new(OkexConn)
	okc.WsConn = ws
	okc.Health = true
	okc.HealthTime = time.Now().Unix()
	return okc
}

type OkexApi struct {
	sync.RWMutex
	//TODO use one connection
	upstreamConns map[string]*OkexConn
}

func NewOkexApi() *OkexApi {
	oa := new(OkexApi)
	oa.upstreamConns = make(map[string]*OkexConn)
	return oa
}

func (this *OkexApi) Start(contract string, table string) error {
	//check conn existence
	key := this.getKey(contract, table)
	if _, ok := this.checkExt(key); ok {
		log.WARNF("[okexapi_start]duplicate sub req for %s-%s", contract, table)
		return nil
	}
	wsConn, err := websocket.Dial(OKEX_SERVER, PROTO_WSS, OKEX_URI)
	if err != nil {
		log.ERRORF("[okexapi_start]connect to okex api error: %v", err)
		return err
	}
	okexConn := NewOkexConn(wsConn)
	this.Lock()
	this.upstreamConns[key] = okexConn
	this.Unlock()
	log.DEBUG(this.upstreamConns)
	//start sub and health check
	go this.Sub(contract, table, okexConn)
	return nil
}

func (this *OkexApi) Sub(contract string, table string, okexConn *OkexConn) {
	ticker := time.NewTicker(HEALTH_CHECK_TIME)
	defer ticker.Stop()
	var (
		key   = this.getKey(contract, table)
		count = 0
	)
	//sub process
	switch table {
	case "orderbook":
		go this.SendFutureUsed(contract, okexConn.WsConn)
	default:
		log.ERRORF("[okexapi]not found table %s", table)
		return
	}
	for {
		select {
		case <-ticker.C:
			//health check
			if okexConn.Health && time.Now().Unix()-okexConn.HealthTime <= HEALTH_GAP {
				go this.SendPingPong(okexConn.WsConn)
				break
			} else {
				log.WARNF("[okexapi]conn %s-%s is in abnormal", contract, table)
				count += 1
			}
			if count > HEALTH_CHECK_RETRY {
				//reconnect
				var (
					wsConn *websocket.Conn
					err    error
				)
			ALWAYSRETRY:
				{
					wsConn, err = websocket.Dial(OKEX_SERVER, PROTO_WSS, OKEX_URI)
					if err != nil {
						log.ERRORF("[okexapi]reconnect error: %v", err)
						goto ALWAYSRETRY
					}
				}
				okexConn = NewOkexConn(wsConn)
				this.Lock()
				this.upstreamConns[key] = okexConn
				this.Unlock()
				count = 0
			}
		default:
			//data channel
			var (
				buf     []byte
				ppRes   PingPongResponse
				dataRes []DataResponse
			)
			err := websocket.Message.Receive(okexConn.WsConn, &buf)
			if err != nil {
				log.ERRORF("[okexapi_sub]receive for %s-%s error: %v", contract, table, err)
				break
			}
			//check event
			json.Unmarshal(buf, &ppRes)
			if ppRes.Event == "pong" {
				okexConn.Health = true
				okexConn.HealthTime = time.Now().Unix()
				count = 0
				break
			}
			//check data
			//TODO channel(table) route
			json.Unmarshal(buf, &dataRes)
			if len(dataRes) == 0 {
				log.ERRORF("[okexapi]receive data empty")
				break
			}
			log.DEBUG(dataRes)
			if okexManager, ok := conn.GConn.Load(key); ok {
				cliLst := okexManager.(conn.ConnManager).DumpConns(contract, table)
				if len(cliLst) == 0 {
					this.Lock()
					delete(this.upstreamConns, key)
					this.Unlock()
					okexConn.WsConn.Close()
				}
				for _, cli := range cliLst {
					go func() {
						_, err := cli.Conn.Write(buf)
						if err != nil {
							log.ERRORF("[okexapi_sub] send data to client error: %v, addr: %s, contract: %s, table: %s", err, cli.RemoteAddr, contract, table)
						}
					}()
				}
			}
		}
	}
}

func (this *OkexApi) getKey(contract string, table string) string {
	return fmt.Sprintf("%s|%s", contract, table)
}

func (this *OkexApi) checkExt(key string) (conn *OkexConn, ok bool) {
	this.RLock()
	defer this.RUnlock()
	conn, ok = this.upstreamConns[key]
	return
}

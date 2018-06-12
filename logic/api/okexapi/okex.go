package okexapi

import (
	"encoding/json"
	"sidekick/tmatrix/logic/api/apitypes"
	"sidekick/tmatrix/logic/conn"
	"strings"
	"sync"
	"time"
	"xframe/log"
	"xframe/server/websocket"
)

type OkexApi struct {
	upstreamConns sync.Map
}

func NewOkexApi() *OkexApi {
	oa := new(OkexApi)
	return oa
}

func (this *OkexApi) Start(contract string, table string, ttype string, depth int, period string) error {
	//check conn existence
	key := this.getKey(contract, table, ttype)
	if ok := this.checkExt(key); ok {
		log.WARNF("[okexapi_start]duplicate sub req for %s-%s-%s", contract, table, ttype)
		return nil
	}
	okexConn := NewOkexConn()
	err := okexConn.Dial(OKEX_SERVER, PROTO_WSS, OKEX_URI)
	if err != nil {
		log.ERRORF("[okexapi_start]connect to okex api error: %v", err)
		return err
	}
	this.upstreamConns.Store(key, okexConn)
	log.DEBUGF("[okexapi]upstream conn: %v", this.upstreamConns)
	//start sub and health check
	go this.Sub(contract, table, ttype, okexConn, depth, period)
	return nil
}

func (this *OkexApi) checkChannel(channel string) string {
	if channel == "ok_sub_futureusd_positions" ||
		channel == "ok_sub_futureusd_trades" ||
		channel == "ok_sub_futureusd_userinfo" {
		return apitypes.API_DATA_LOGIN
	} else if strings.Index(channel, "depth") != -1 {
		return apitypes.API_DATA_ORDERBOOK
	} else if strings.Index(channel, "ticker") != -1 {
		return apitypes.API_DATA_TICKER
	}
	return ""
}

func (this *OkexApi) Sub(contract string, table string, ttype string, okexConn *OkexConn, depth int, period string) {
	ticker := time.NewTicker(apitypes.HEALTH_CHECK_TIME)
	defer ticker.Stop()
	var (
		key string = this.getKey(contract, table, ttype)
		//ppRes         PingPongResponse
		dataRes       []DataResponse
		dataCommonRes []DataCommonRes
		buf           []byte
		pushBuf       []byte
	)
	//sub process
	switch table {
	case apitypes.API_DATA_ORDERBOOK:
		go this.SendFutureUsd(contract, ttype, okexConn.WsConn, depth)
	case apitypes.API_DATA_LOGIN:
		go this.SendLogin(okexConn.WsConn)
	case apitypes.API_DATA_TICKER:
		go this.SendFutureUsdTicker(contract, ttype, okexConn.WsConn)
	default:
		log.ERRORF("[okexapi]not found table %s", table)
		return
	}
	for {
		select {
		case <-ticker.C:
			//FIXME after adding health check, we can not receive data any more
			log.DEBUGF("[okexapi]start to check connection")
			break
			//health check
			/*if okexConn.CheckHealth() {
				go this.SendPingPong(okexConn.WsConn)
				break
			} else {
				log.WARNF("[okexapi]conn %s-%s-%s is in abnormal", contract, table, ttype)
				okexConn.SetUnhealthy()
			}
			if okexConn.UnHealthCount() && !okexConn.IsRetry() {
				//reconnect
				log.DEBUGF("[okexapi]reconnect to okex")
				okexConn.SetRetry()
				var (
					newOkexConn = NewOkexConn()
					err         error
				)
			ALWAYSRETRY:
				{
					err = newOkexConn.Dial(OKEX_SERVER, PROTO_WSS, OKEX_URI)
					if err != nil {
						log.ERRORF("[okexapi]reconnect error: %v", err)
						goto ALWAYSRETRY
					}
				}
				go okexConn.Close()
				this.upstreamConns.Store(key, newOkexConn)
				//recursion
				go this.Sub(contract, table, ttype, newOkexConn)
				return
			}*/
		default:
			//data channel
			err := websocket.Message.Receive(okexConn.WsConn, &buf)
			if err != nil {
				log.ERRORF("[okexapi_sub]receive for %s-%s-%s error: %v", contract, table, ttype, err)
				var (
					newOkexConn = NewOkexConn()
					err         error
				)
			ALWAYSRETRY:
				{
					err = newOkexConn.Dial(OKEX_SERVER, PROTO_WSS, OKEX_URI)
					if err != nil {
						log.ERRORF("[okexapi]reconnect error: %v", err)
						goto ALWAYSRETRY
					}
				}
				go okexConn.Close()
				this.upstreamConns.Store(key, newOkexConn)
				//recursion
				go this.Sub(contract, table, ttype, newOkexConn, depth, period)
				return
			}
			//check event
			/*json.Unmarshal(buf, &ppRes)
			if ppRes.Event == "pong" {
				log.DEBUGF("[okexapi]receive ping-pong response")
				okexConn.SetHealthy()
				break
			}*/
			//check channel
			json.Unmarshal(buf, &dataCommonRes)
			if len(dataCommonRes) == 0 {
				log.ERRORF("[okexapi]receive data empty")
				break
			}
			log.DEBUG("[okexapi]data from okex:", dataCommonRes)
			if okexManager, ok := conn.GConn.Load("okex"); ok {
				log.DEBUG("[okexapi]start to dump okex client list")
				cliLst := okexManager.(conn.ConnManager).DumpConns(contract, table, ttype)
				if len(cliLst) == 0 {
					this.upstreamConns.Delete(key)
					okexConn.Close()
					return
				}
				//TODO register callback for different push data
				//generate push data
				switch this.checkChannel(dataCommonRes[0].Channel) {
				case apitypes.API_DATA_LOGIN:
					//json.Unmarshal(buf, &dataCommonRes)
					pushTable := this._getLoginTable(dataCommonRes[0].Channel)
					pushBuf, err = this.generateLoginPushData(dataCommonRes, "okex", contract, pushTable, ttype)
					if err != nil {
						log.ERRORF("[okexapi]generate login push data error: %v", err)
						break
					}
				case apitypes.API_DATA_ORDERBOOK:
					json.Unmarshal(buf, &dataRes)
					pushBuf, err = this.generateQuotePushData(dataRes, "okex", contract, table, ttype, depth)
					if err != nil {
						log.ERRORF("[okexapi]generate quote push data error: %v", err)
						break
					}
				case apitypes.API_DATA_TICKER:
					//json.Unmarshal(buf, &dataCommonRes)
					pushBuf, err = this.generateTickerPushData(dataCommonRes, "okex", contract, table, ttype)
					if err != nil {
						log.ERRORF("[okexapi]generate ticker push data error: %v", err)
					}
				default:
					log.ERRORF("[okexapi]channel %s not found", dataCommonRes[0].Channel)
					break
				}
				//push to clients
				for _, cli := range cliLst {
					log.DEBUGF("[okexapi]send to client: %v", cli)
					go func(buf []byte) {
						_, err := cli.Conn.Write(buf)
						if err != nil {
							log.ERRORF("[okexapi_sub] send data to client error: %v, addr: %s, contract: %s, table: %s, type: %s, depth: %d, period: %s", err, cli.RemoteAddr, contract, table, ttype, depth, period)
						}
					}(pushBuf)
				}
			}
		}
	}
}

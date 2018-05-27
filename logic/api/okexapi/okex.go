package okexapi

import (
	"encoding/json"
	"fmt"
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

func (this *OkexApi) Start(contract string, table string, ttype string) error {
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
	go this.Sub(contract, table, ttype, okexConn)
	return nil
}

func (this *OkexApi) generateQuotePushData(dataRes []DataResponse, market string, contract string, table string, ttype string) ([]byte, error) {
	var (
		dataPush DataPush
	)
	dataPush.Msg = "quote"
	dataPush.Market = market
	dataPush.Table = table
	dataPush.Data = make([]QuotePush, 0)
	for _, data := range dataRes {
		var tmpQuotePush QuotePush
		tmpQuotePush.Ts = data.Data.Ts
		tmpQuotePush.Asks = data.Data.Asks
		tmpQuotePush.Bids = data.Data.Bids
		tmpQuotePush.Contract = contract
		tmpQuotePush.Type = ttype
		dataPush.Data = append(dataPush.Data, tmpQuotePush)
	}
	buf, err := json.Marshal(dataPush)
	return buf, err
}

func (this *OkexApi) generateLoginPushData(dataCommonRes []DataCommonRes, market string, contract string, table string, ttype string) ([]byte, error) {
	var (
		dataCommonPush DataCommonPush
	)
	dataCommonPush.Msg = "login"
	dataCommonPush.Table = table
	for _, data := range dataCommonRes {
		dataCommonPush.Data = append(dataCommonPush.Data, data.Data)
	}
	buf, err := json.Marshal(dataCommonPush)
	return buf, err
}

func (this *OkexApi) checkChannel(channel string) string {
	if channel == "ok_sub_futureusd_positions" ||
		channel == "ok_sub_futureusd_trades" ||
		channel == "ok_sub_futureusd_userinfo" {
		return "login"
	} else if strings.Index(channel, "depth") != -1 {
		return "depth"
	}
	return ""
}

func (this *OkexApi) Sub(contract string, table string, ttype string, okexConn *OkexConn) {
	ticker := time.NewTicker(HEALTH_CHECK_TIME)
	defer ticker.Stop()
	var (
		key string = this.getKey(contract, table, ttype)
		//ppRes         PingPongResponse
		dataRes       []DataResponse
		dataCommonRes []DataCommonRes
		channelRes    []ChannelRes
		buf           []byte
		pushBuf       []byte
	)
	//sub process
	switch table {
	case apitypes.API_DATA_ORDERBOOK:
		go this.SendFutureUsed(contract, ttype, okexConn.WsConn)
	case apitypes.API_DATA_LOGIN:
		go this.SendLogin(okexConn.WsConn)
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
				go this.Sub(contract, table, ttype, newOkexConn)
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
			json.Unmarshal(buf, &channelRes)
			if len(channelRes) == 0 {
				log.ERRORF("[okexapi]receive data empty")
				break
			}
			log.DEBUG("[okexapi]data from okex:", channelRes)
			if okexManager, ok := conn.GConn.Load("okex"); ok {
				log.DEBUG("[okexapi]start to dump okex client list")
				cliLst := okexManager.(conn.ConnManager).DumpConns(contract, table, ttype)
				if len(cliLst) == 0 {
					this.upstreamConns.Delete(key)
					okexConn.Close()
					return
				}
				//generate push data
				if this.checkChannel(dataCommonRes[0].Channel) == "login" {
					json.Unmarshal(buf, &dataCommonRes)
					pushBuf, err = this.generateLoginPushData(dataCommonRes, "okex", contract, table, ttype)
					if err != nil {
						log.ERRORF("[okexapi]generate login push data error: %v", err)
						break
					}
				} else if this.checkChannel(dataCommonRes[0].Channel) == "depth" {
					json.Unmarshal(buf, &dataRes)
					pushBuf, err = this.generateQuotePushData(dataRes, "okex", contract, table, ttype)
					if err != nil {
						log.ERRORF("[okexapi]generate quote push data error: %v", err)
						break
					}
				} else {
					log.ERRORF("[okexapi]channel %s not found", dataCommonRes[0].Channel)
					break
				}
				//push to clients
				for _, cli := range cliLst {
					log.DEBUGF("[okexapi]send to client: %v", cli)
					go func(buf []byte) {
						_, err := cli.Conn.Write(buf)
						if err != nil {
							log.ERRORF("[okexapi_sub] send data to client error: %v, addr: %s, contract: %s, table: %s, type: %s", err, cli.RemoteAddr, contract, table, ttype)
						}
					}(pushBuf)
				}
			}
		}
	}
}

func (this *OkexApi) getKey(contract string, table string, ttype string) string {
	return fmt.Sprintf("%s|%s|%s", contract, table, ttype)
}

func (this *OkexApi) checkExt(key string) (ok bool) {
	_, ok = this.upstreamConns.Load(key)
	return
}

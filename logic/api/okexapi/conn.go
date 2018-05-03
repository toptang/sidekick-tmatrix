package okexapi

import (
	"time"
	"xframe/server/websocket"
)

type OkexConn struct {
	WsConn      *websocket.Conn
	Health      bool
	HealthTime  int64
	HealthCount int
	Retry       bool
}

func NewOkexConn() *OkexConn {
	okc := new(OkexConn)
	okc.Health = true
	okc.HealthTime = time.Now().Unix()
	return okc
}

func (this *OkexConn) Dial(server string, scheme string, uri string) error {
	wsConn, err := websocket.Dial(OKEX_SERVER, PROTO_WSS, OKEX_URI)
	if err != nil {
		return err
	}
	this.WsConn = wsConn
	return nil
}

func (this *OkexConn) SetHealthy() {
	this.Health = true
	this.HealthTime = time.Now().Unix()
	this.HealthCount = 0
}

func (this *OkexConn) CheckHealth() bool {
	return this.Health && time.Now().Unix()-this.HealthTime <= HEALTH_GAP
}

func (this *OkexConn) SetUnhealthy() {
	this.HealthCount += 1
}

func (this *OkexConn) UnHealthCount() bool {
	return this.HealthCount > HEALTH_CHECK_RETRY
}

func (this *OkexConn) SetRetry() {
	this.Retry = true
}

func (this *OkexConn) IsRetry() bool {
	return this.Retry
}

func (this *OkexConn) Close() {
	this.WsConn.Close()
}

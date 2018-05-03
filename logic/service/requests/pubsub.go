package requests

import "xframe/log"

type PubSubReq struct {
	Msg  string `json:"msg"`
	Uuid string `json:"uuid"`
	Data struct {
		Market string `json:"market"`
		Type   string `json:"type"`
		Symbol string `json:"symbol"` //contract
		Table  string `json:"table"`
	} `json:"data"`
}

func (this PubSubReq) CheckRouter() bool {
	return this.Msg == ""
}

func (this PubSubReq) CheckParams() bool {
	return (this.Data.Market == "" ||
		this.Data.Symbol == "" ||
		this.Data.Table == "" ||
		this.Data.Type == "")
}

func (this PubSubReq) Dump() {
	log.DEBUGF("[pubsub_req] Msg: %s, UUID: %s, Market: %s, Contract: %s, Table: %s",
		this.Msg,
		this.Uuid,
		this.Data.Market,
		this.Data.Symbol,
		this.Data.Table)
}

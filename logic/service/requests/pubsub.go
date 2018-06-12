package requests

import "xframe/log"

type Option struct {
	Depth  int    `json:"depth"`
	Period string `json:"period"`
}

type PubSubReq struct {
	Msg  string `json:"msg"`  //action
	Uuid string `json:"uuid"` //request id
	Data struct {
		Market   string `json:"market"` //交易所
		Type     string `json:"type"`   //货品: 现货/期货
		Symbol   string `json:"symbol"` //合约
		Table    string `json:"table"`  //数据类型
		Optional Option `json:"optional"`
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

package requests

type PubSubReq struct {
	Msg  string `json:"msg"`
	Uuid string `json:"uuid"`
	Data struct {
		Market string `json:"market"`
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
		this.Data.Table == "")
}

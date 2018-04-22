package requests

type PubSubReq struct {
	Msg  string `json:"msg"`
	Data struct {
		Market   string `json:"market"`
		Contract string `json:"contract"`
	} `json:"data"`
}

func (this PubSubReq) CheckRouter() bool {
	return this.Msg == ""
}

func (this PubSubReq) CheckParams() bool {
	return (this.Data.Market == "" ||
		this.Data.Contract == "")
}

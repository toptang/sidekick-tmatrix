package response

import (
	"encoding/json"
	"xframe/server/websocket"
)

var (
	ERROR_MSG = "error"
)

type BaseResponse struct {
	Msg     string `json:"msg"`
	RetCode int    `json:"retcode"`
	Message string `json:"message"`
}

//TODO: 上游数据格式
type Pricing struct {
}

type DataResponse struct {
	Msg  string `json:"msg"`
	Data struct {
		Market string `json:"market"`
		Pricing
	} `json:"data"`
}

//-------------------------

func formatBaseResponse(route string, retcode int) (res []byte, err error) {
	var message string
	if retcode != err.SUCCESS {
		message = err.ErrMap[retcode]
	}
	var baseResponse = BaseResponse{
		Msg:     route,
		RetCode: retcode,
		Message: message,
	}
	res, err = json.Marshal(baseResponse)
	return
}

func DoBaseResponse(route string, retcode int, ws *websocket.Conn) {
	res, err := formatBaseResponse(route, retcode)
	if err != nil {
		log.ERRORF("[response]send response error: %v, retcode: %d", err, route)
		return
	}
	ws.Write(res)
}

//---------------------------

func DoDataResponse(buf []byte, ws *websocket.Conn) {
	ws.Write(buf)
}

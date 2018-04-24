package response

import (
	"encoding/json"
	"sidekick/tmatrix/logic/service/svcerr"
	"xframe/log"
	"xframe/server/websocket"
)

type BaseResponse struct {
	Msg  string `json:"msg"`
	Err  string `json:"err"`
	Uuid string `json:"uuid"`
}

//NEVER USED
type DataResponse struct {
	Msg string `json:"msg"`
}

//-------------------------

func formatBaseResponse(route string, retcode int, uuid string) (res []byte, err error) {
	var message string
	if retcode != svcerr.SUCCESS {
		message = svcerr.ErrMap[retcode]
	}
	var baseResponse = BaseResponse{
		Msg:  route,
		Err:  message,
		Uuid: uuid,
	}
	res, err = json.Marshal(baseResponse)
	return
}

func DoBaseResponse(route string, retcode int, uuid string, ws *websocket.Conn) {
	res, err := formatBaseResponse(route, retcode, uuid)
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

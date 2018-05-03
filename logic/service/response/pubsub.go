package response

import (
	"encoding/json"
	"sidekick/tmatrix/logic/service/svcerr"
	"xframe/log"
	"xframe/server/websocket"
)

type ErrS struct {
	Code   int    `json:"code"`
	ErrMsg string `json:"what"`
}

type BaseResponse struct {
	Msg  string `json:"msg"`
	Uuid string `json:"uuid"`
	Err  ErrS   `json:"err"`
}

//-------------------------

func formatBaseResponse(route string, retcode int, uuid string) (res []byte, err error) {
	var (
		message string
		rspCode int = retcode
	)
	if retcode != svcerr.SUCCESS {
		rspCode = -1 * retcode
		message = svcerr.ErrMap[retcode]
	}
	errS := ErrS{
		Code:   rspCode,
		ErrMsg: message,
	}
	var baseResponse = BaseResponse{
		Msg:  "rsp" + route,
		Uuid: uuid,
		Err:  errS,
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

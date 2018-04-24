package svcerr

const (
	SUCCESS = iota

	ROUTE_ERROR          = -1000
	INTERNAL_ERROR       = -1001
	INPUT_ERROR          = -1002
	CONN_NOT_FOUND_ERROR = -1003
	API_NOT_FOUND_ERROR  = -1004
)

var (
	ErrMap = map[int]string{
		ROUTE_ERROR:          "protocol route not found",
		INTERNAL_ERROR:       "internal error",
		INPUT_ERROR:          "params error",
		CONN_NOT_FOUND_ERROR: "market connection manager not found",
		API_NOT_FOUND_ERROR:  "upstream api manager not found",
	}
)

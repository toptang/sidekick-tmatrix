package svcerr

const (
	SUCCESS = iota

	ROUTE_ERROR = iota + 1000
	INTERNAL_ERROR
	INPUT_ERROR
	CONN_NOT_FOUND_ERROR
	API_NOT_FOUND_ERROR
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

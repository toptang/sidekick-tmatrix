package err

var (
	SUCCESS = iota

	ROUTE_ERROR    = -1000
	INTERNAL_ERROR = -1001
	INPUT_ERROR    = -1002
)

var (
	ErrMap = map[int]string{
		ROUTE_ERROR:    "protocol route not found",
		INTERNAL_ERROR: "internal error",
		INPUT_ERROR:    "params error",
	}
)

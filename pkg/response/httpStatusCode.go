package response

import "net/http"

const (
	ErrCodeSuccess      = 20001
	ErrCodeParamInvalid = 20003
	ErrCodeUnauthorized = 20004
	ErrCodeNotFound     = 20005
	ErrCodeServerError  = 20006
)

var msg = map[int]string{
	ErrCodeSuccess:      "Success",
	ErrCodeParamInvalid: "Parameter is invalid",
	ErrCodeUnauthorized: "Unauthorized",
	ErrCodeNotFound:     "Resource not found",
	ErrCodeServerError:  "Internal server error",
}

var errCodeToHTTPStatus = map[int]int{
	ErrCodeSuccess:      http.StatusOK,
	ErrCodeParamInvalid: http.StatusBadRequest,
	ErrCodeUnauthorized: http.StatusUnauthorized,
	ErrCodeNotFound:     http.StatusNotFound,
	ErrCodeServerError:  http.StatusInternalServerError,
}

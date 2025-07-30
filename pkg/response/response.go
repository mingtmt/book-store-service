package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ResponseData{
		Code:    ErrCodeSuccess,
		Message: msg[ErrCodeSuccess],
		Data:    data,
	})
}

func SuccessWithCode(c *gin.Context, code int, data interface{}) {
	c.JSON(http.StatusOK, ResponseData{
		Code:    code,
		Message: getMsg(code),
		Data:    data,
	})
}

func Error(c *gin.Context, code int, customMsg ...string) {
	message := getMsg(code)
	if len(customMsg) > 0 && customMsg[0] != "" {
		message = customMsg[0]
	}
	httpStatus := HTTPStatusFromErrCode(code)
	c.JSON(httpStatus, ResponseData{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func getMsg(code int) string {
	if m, ok := msg[code]; ok {
		return m
	}
	if code == ErrCodeSuccess {
		return "Success"
	}
	return "An error occurred"
}

func HTTPStatusFromErrCode(code int) int {
	if status, ok := errCodeToHTTPStatus[code]; ok {
		return status
	}
	return http.StatusInternalServerError
}

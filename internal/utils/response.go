package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorCode string

const (
	ErrCodeBadRequest ErrorCode = "BAD_REQUEST"
	ErrCodeNotFound   ErrorCode = "NOT_FOUND"
	ErrCodeConflict   ErrorCode = "CONFLICT"
	ErrCodeInternal   ErrorCode = "INTERNAL_SERVER_ERR"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (ae *AppError) Error() string {
	return ""
}

func NewError(code ErrorCode, message string) error {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func WrapError(code ErrorCode, message string, err error) error {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func ResponseError(ctx *gin.Context, err error) {
	if appErr, oke := err.(*AppError); oke {
		status := httpStatusFromCode(appErr.Code)
		response := gin.H{
			"error": appErr.Message,
			"code":  appErr.Code,
		}

		if appErr.Err != nil {
			response["detail"] = appErr.Err.Error()
		}

		ctx.JSON(status, response)
		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
		"code":  ErrCodeInternal,
	})
}

func ResponseSuccess(ctx *gin.Context, status int, data any) {
	ctx.JSON(status, gin.H{
		"status": "success",
		"data":   data,
	})
}

func ResponseStatusCode(ctx *gin.Context, status int) {
	ctx.Status(status)
}

func ResponseValidator(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusBadRequest, data)
}

func httpStatusFromCode(code ErrorCode) int {
	switch code {
	case ErrCodeBadRequest:
		return http.StatusBadRequest
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	ResponseStatus ResponseStatusType `json:"responseStatus"`
	Error          *ErrorType         `json:"error,omitempty"`
}

type ResponseStatusType string

type BaseDataResponse struct {
	BaseResponse
	Data interface{} `json:"data"`
}

const (
	RespStatus_SUCCESS ResponseStatusType = "SUCCESS"
	RespStatus_ERROR   ResponseStatusType = "ERROR"
	RespStatus_EXPIRED ResponseStatusType = "EXPIRED"
)

type ErrorType struct {
	Code        ErrorCodeType `json:"code"`
	Description string        `json:"description"`
	Errors      *[]ErrorType  `json:"errors,omitempty"`
}

type BaseSuccessResponse struct {
	BaseResponse
	Data interface{} `json:"data,omitempty"`
}

func SendSuccessResponse(context *gin.Context, data interface{}) {
	context.SecureJSON(http.StatusOK, BaseSuccessResponse{
		BaseResponse: BaseResponse{
			ResponseStatus: RespStatus_SUCCESS,
		},
		Data: data,
	})
}

func SendErrorResponse(context *gin.Context, code int, errorCode ErrorCodeType, err error) {
	context.JSON(code, BaseSuccessResponse{
		BaseResponse: BaseResponse{
			ResponseStatus: RespStatus_ERROR,
			Error: &ErrorType{
				Code:        errorCode,
				Description: err.Error(),
			},
		},
	})
}

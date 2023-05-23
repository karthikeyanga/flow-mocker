package service

import (
	"mocker/common"
	"mocker/config"
	"mocker/core"
	"mocker/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewApiAction(context *gin.Context) {
	ac := config.GetAppContext(context)
	req := ApiRequest{}
	if !util.BindRequestForApiResponse(context, &req) {
		return
	}
	api, err := core.NewApi(ac, req.Route, req.Method, req.Status, req.ResponseHeaders, req.ResponseBody)
	if err != nil {
		util.SetErrorResponseFromError(context, err, "new api creation")
		return
	}
	if api == nil {
		context.AbortWithStatus(http.StatusNotFound)
		return
	}
	util.SetTypedApiSuccessResponse(context, api)
}

func UpdateApiAction(context *gin.Context) {
	ac := config.GetAppContext(context)

	apiId := context.Param("apiId")
	api, err := core.GetApi(ac, common.StringToInt64(apiId))
	if err != nil {
		util.SetErrorResponseFromError(context, err, "getting api for updation")
		return
	}
	if api == nil {
		context.AbortWithStatus(http.StatusNotFound)
		return
	}
	req := ApiRequest{}
	if !util.BindRequestForApiResponse(context, &req) {
		return
	}

	if !api.IsEnabled() {
		context.JSON(http.StatusInternalServerError, common.BaseDataResponse{
			BaseResponse: common.BaseResponse{
				ResponseStatus: common.RespStatus_ERROR,
				Error: &common.ErrorType{
					Code:        "ERROR",
					Description: "not enabled - " + apiId,
				},
			},
			Data: nil,
		})
		return
	}
	err = api.Update(ac, req.Route, req.Method, req.Status, req.ResponseHeaders, req.ResponseBody)
	if err != nil {
		util.SetErrorResponseFromError(context, err, "updating api")
		return
	}
	util.SetTypedApiSuccessResponse(context, api)
}

func GetApiAction(context *gin.Context) {
	ac := config.GetAppContext(context)
	apiId := context.Param("apiId")
	api, err := core.GetApi(ac, common.StringToInt64(apiId))
	if err != nil {
		util.SetErrorResponseFromError(context, err, "getting api")
		return
	}
	if api == nil {
		context.AbortWithStatus(http.StatusNotFound)
		return
	}
	if !api.IsEnabled() {
		util.SetTypedApiSuccessResponse(context, nil)
		return
	}
	util.SetTypedApiSuccessResponse(context, api)
}

func EnableApiAction(context *gin.Context) {
	ac := config.GetAppContext(context)
	apiId := context.Param("apiId")
	api, err := core.GetApi(ac, common.StringToInt64(apiId))
	if err != nil {
		util.SetErrorResponseFromError(context, err, "getting api for enabling/disabling")
		return
	}
	if api == nil {
		context.AbortWithStatus(http.StatusNotFound)
		return
	}
	action := context.Param("action")
	if api.IsEnabled() != (action == "enable") {
		switch action {
		case "enable":
			err = api.Enable(ac)
		case "disable":
			err = api.Disable(ac)
		}
		if err != nil {
			util.SetErrorResponseFromError(context, err, "disabling api")
			return
		}
	}
	util.SetTypedApiSuccessResponse(context, nil)
}

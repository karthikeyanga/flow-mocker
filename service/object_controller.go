package service

import (
	"mocker/common"
	"mocker/config"
	"mocker/core"
	"mocker/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewObjectAction(context *gin.Context) {
	ac := config.GetAppContext(context)
	req := common.Json{}
	if !util.BindRequestForApiResponse(context, &req) {
		return
	}
	obj, err := core.NewObject(ac, &req)
	if err != nil {
		util.SetErrorResponseFromError(context, err, "new object creation")
		return
	}
	if obj == nil {
		context.AbortWithStatus(http.StatusNotFound)
		return
	}
	util.SetTypedApiSuccessResponse(context, obj)

}

func UpdateObjectAction(context *gin.Context) {
	ac := config.GetAppContext(context)
	req := common.Json{}
	objId := context.Param("objectId")
	if !util.BindRequestForApiResponse(context, &req) {
		return
	}

	obj, err := core.GetObject(ac, common.StringToInt64(objId))
	if err != nil {
		util.SetErrorResponseFromError(context, err, "getting object for updation")
		return
	}
	if obj == nil {
		context.AbortWithStatus(http.StatusNotFound)
		return
	}
	if !obj.IsEnabled() {
		context.JSON(http.StatusInternalServerError, common.BaseDataResponse{
			BaseResponse: common.BaseResponse{
				ResponseStatus: common.RespStatus_ERROR,
				Error: &common.ErrorType{
					Code:        "ERROR",
					Description: "not enabled - " + objId,
				},
			},
			Data: nil,
		})
		return
	}
	err = obj.Update(ac, &req)
	if err != nil {
		util.SetErrorResponseFromError(context, err, "updating object")
		return
	}
	util.SetTypedApiSuccessResponse(context, obj)
}

func GetObjectAction(context *gin.Context) {
	ac := config.GetAppContext(context)
	objId := context.Param("objectId")
	obj, err := core.GetObject(ac, common.StringToInt64(objId))
	if err != nil {
		util.SetErrorResponseFromError(context, err, "getting object")
		return
	}
	if obj == nil {
		context.AbortWithStatus(http.StatusNotFound)
		return
	}
	if !obj.IsEnabled() {
		util.SetTypedApiSuccessResponse(context, nil)
		return
	}
	util.SetTypedApiSuccessResponse(context, obj)
}

func EnableObjectAction(context *gin.Context) {
	ac := config.GetAppContext(context)
	objId := context.Param("objectId")
	obj, err := core.GetObject(ac, common.StringToInt64(objId))
	if err != nil {
		util.SetErrorResponseFromError(context, err, "getting object for disabling")
		return
	}
	if obj == nil {
		context.AbortWithStatus(http.StatusNotFound)
		return
	}
	action := context.Param("action")
	if obj.IsEnabled() != (action == "enable") {
		switch action {
		case "enable":
			err = obj.Enable(ac)
		case "disable":
			err = obj.Disable(ac)
		}
		if err != nil {
			util.SetErrorResponseFromError(context, err, "disabling object")
			return
		}
	}
	util.SetTypedApiSuccessResponse(context, nil)
}

func ShowLoadedObjectsAction(context *gin.Context) {
	ac := config.GetAppContext(context)
	objects := ac.GetObjects()
	util.SetTypedApiSuccessResponse(context, objects)
}

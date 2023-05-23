package service

import (
	"io/ioutil"
	"mocker/common"
	"mocker/config"
	"mocker/core"
	"mocker/util"
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/gin-gonic/gin"
)

func LoadNewFlowAction(context *gin.Context) {
	ac := config.GetAppContext(context)
	req := FlowRequest{}
	switch context.GetHeader("Content-Type") {
	case "text/yaml", "text/x-yaml", "application/x-yaml", "text/yml", "text/x-yml", "application/x-yml":
		defer context.Request.Body.Close()
		b, err := ioutil.ReadAll(context.Request.Body)
		if err != nil {
			util.SetErrorResponseFromError(context, err, "reading yaml for new flow ")
			return
		}
		err = yaml.Unmarshal(b, &req)
		if err != nil {
			util.SetErrorResponseFromError(context, err, "reading yaml for new flow ")
			return
		}
	default:
		if !util.BindRequestForApiResponse(context, &req) {
			return
		}
	}

	f, err := core.NewFlow(ac, req.Title, req.Identifier, req.ObjectIds(), req.Config)
	if err != nil {
		util.SetErrorResponseFromError(context, err, "creating new flow for loading")
		return
	}
	if f == nil {
		context.AbortWithStatus(http.StatusNotFound)
		return
	}
	err = f.Load(ac)
	if err != nil {
		util.SetErrorResponseFromError(context, err, "loading new flow ")
		return
	}
	util.SetTypedApiSuccessResponse(context, f.Id())
}

func UpdateFlowAction(context *gin.Context) {
	ac := config.GetAppContext(context)
	flowId := context.Param("flowId")
	f, err := core.GetFlow(ac, common.StringToInt64(flowId))
	if err != nil {
		util.SetErrorResponseFromError(context, err, "getting flow for loading")
		return
	}
	if f == nil {
		context.AbortWithStatus(http.StatusNotFound)
		return
	}
	req := FlowRequest{}
	contentType := context.GetHeader("Content-Type")
	switch contentType {
	case "text/yaml", "text/x-yaml", "application/x-yaml", "text/yml", "text/x-yml", "application/x-yml":
		defer context.Request.Body.Close()
		b, err := ioutil.ReadAll(context.Request.Body)
		if err != nil {
			util.SetErrorResponseFromError(context, err, "reading yaml for new flow ")
			return
		}
		err = yaml.Unmarshal(b, &req)
		if err != nil {
			util.SetErrorResponseFromError(context, err, "reading yaml for new flow ")
			return
		}
	default:
		if !util.BindRequestForApiResponse(context, &req) {
			return
		}
	}

	err = f.Update(ac, req.Title, req.Identifier, req.ObjectIds(), req.Config)
	if err != nil {
		util.SetErrorResponseFromError(context, err, "loading new flow ")
		return
	}
	util.SetTypedApiSuccessResponse(context, f.Id())
}

func GetFlowAction(context *gin.Context) {
	ac := config.GetAppContext(context)
	flowId := context.Param("flowId")
	f, err := core.GetFlow(ac, common.StringToInt64(flowId))
	if err != nil {
		util.SetErrorResponseFromError(context, err, "getting flow for loading")
		return
	}
	if f == nil {
		context.AbortWithStatus(http.StatusNotFound)
		return
	}

	util.SetTypedApiSuccessResponse(context, f)
}

func LoadFlowAction(context *gin.Context) {
	ac := config.GetAppContext(context)
	flowId := context.Param("flowId")
	f, err := core.GetFlow(ac, common.StringToInt64(flowId))
	if err != nil {
		util.SetErrorResponseFromError(context, err, "getting flow for loading")
		return
	}
	if f == nil {
		context.AbortWithStatus(http.StatusNotFound)
		return
	}
	err = f.Load(ac)
	if err != nil {
		util.SetErrorResponseFromError(context, err, "loading flow")
		return
	}
	util.SetTypedApiSuccessResponse(context, f.Id())
}

func UnloadFlowAction(context *gin.Context) {
	ac := config.GetAppContext(context)
	flowId := context.Param("flowId")
	f, err := core.GetFlow(ac, common.StringToInt64(flowId))
	if err != nil {
		util.SetErrorResponseFromError(context, err, "getting flow for unloading")
		return
	}
	if f == nil {
		context.AbortWithStatus(http.StatusNotFound)
		return
	}
	err = f.Unload(ac)
	if err != nil {
		util.SetErrorResponseFromError(context, err, "unloading flow")
		return
	}
	util.SetTypedApiSuccessResponse(context, nil)
}

func ShowLoadedFlowsAction(context *gin.Context) {
	ac := config.GetAppContext(context)
	flows := ac.GetFlows()
	util.SetTypedApiSuccessResponse(context, flows)
}

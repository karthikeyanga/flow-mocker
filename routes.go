package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"mocker/common"
	"mocker/config"
	"mocker/core"
	"mocker/service"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func RouterPatternsInit(ac *config.AppContext, router *gin.Engine) {
	log.Info("In Routes")

	//Internal APIs
	apiRouter := router.Group("/api") // /api
	apisV1 := apiRouter.Group("/v1")  // /api/v1
	{
		apisV1.GET("/ping", func(context *gin.Context) {
			context.String(http.StatusOK, "pong")
		})

	}
	objectApis := apisV1.Group("/objects")
	{
		objectApis.POST("", service.NewObjectAction)
		objectApis.GET("", service.ShowLoadedObjectsAction)
		objectApis.GET("/:objectId", service.GetObjectAction)
		objectApis.POST("/:objectId", service.UpdateObjectAction)
		objectApis.POST("/:objectId/:action", service.EnableObjectAction)
	}
	flowApis := apisV1.Group("/flows")
	{
		flowApis.POST("", service.LoadNewFlowAction)
		flowApis.GET("", service.ShowLoadedFlowsAction)
		flowApis.GET("/:flowId", service.GetFlowAction)
		flowApis.POST("/:flowId", service.UpdateFlowAction)
		flowApis.POST("/:flowId/load", service.LoadFlowAction)
		flowApis.POST("/:flowId/unload", service.UnloadFlowAction)
	}
	apiApis := apisV1.Group("/apis")
	{
		apiApis.POST("", service.NewApiAction)
		apiApis.GET("/:apiId", service.GetApiAction)
		apiApis.POST("/:apiId", service.UpdateApiAction)
		apiApis.POST("/:apiId/:action", service.EnableApiAction)
	}
	apisV1.POST("/whiteboard", service.WhiteboardAction)
	mocksApis := apisV1.Group("/mock")
	for id, api := range ac.Apis {
		switch api.Method() {
		case "GET":
			mocksApis.GET(api.Route(), mockHandler(ac, id))
		case "POST":
			mocksApis.POST(api.Route(), mockHandler(ac, id))
		case "HEAD":
			mocksApis.HEAD(api.Route(), mockHandler(ac, id))
		case "OPTIONS":
			mocksApis.OPTIONS(api.Route(), mockHandler(ac, id))
		case "PATCH":
			mocksApis.PATCH(api.Route(), mockHandler(ac, id))
		case "PUT":
			mocksApis.PUT(api.Route(), mockHandler(ac, id))
		case "DELETE":
			mocksApis.DELETE(api.Route(), mockHandler(ac, id))
		case "Any":
			mocksApis.Any(api.Route(), mockHandler(ac, id))
		}
	}

	// Internal Config Apis

	// partner user login using puid (if linked) >can have a secret

}

func mockHandler(mac *config.AppContext, apiId int64) gin.HandlerFunc {
	api := mac.Apis[apiId]
	return func(context *gin.Context) {
		ac := config.GetAppContext(context)
		//Get the requests
		method := context.Request.Method
		reqHeader := context.Request.Header
		pathParams := map[string]string{}
		for _, p := range context.Params {
			pathParams[p.Key] = p.Value
		}

		getParams := context.Request.URL.Query()
		jsonBody := common.Json{}
		var postForm url.Values
		//Will support post form and json
		//if header of content is there then lets see if its json or else we try to get the post form
		if context.GetHeader("Content-Type") == "application/json" {
			err := context.BindJSON(&jsonBody)
			if err != nil {
				ac.Log.Errorln("error in req json", err)
			}
		} else {
			_ = context.Request.ParseForm()
			postForm = context.Request.PostForm
		}

		data := common.ResponseVariables{
			Method:     method,
			Header:     reqHeader,
			JSON:       jsonBody,
			Query:      getParams,
			PathParams: pathParams,
			Form:       postForm,
		}
		funcMap := core.GetRespRenderFunctions(ac, "RW")
		renderTemplate := func(key, t string, doExecute bool) (*template.Template, string, error) {
			//response
			temp, err := template.New(ac.GetRequestId() + "_" + t).Funcs(funcMap).Parse(t)
			if err != nil {
				return nil, "", err
			}
			if doExecute {
				var buf bytes.Buffer
				err = temp.Execute(&buf, data)
				return temp, buf.String(), err
			}
			return temp, "", nil
		}

		renderErr := func(key, t string, err error) {
			ac.Log.Errorln("error while compiling", key, "template", err)
			context.JSON(http.StatusInternalServerError, common.BaseDataResponse{
				BaseResponse: common.BaseResponse{
					ResponseStatus: common.RespStatus_ERROR,
					Error: &common.ErrorType{
						Code:        "ERROR",
						Description: "error while compiling " + key + " template",
					},
				},
				Data: map[string]interface{}{
					"error":    err,
					"template": t,
				},
			})
			return
		}

		resp := api.ResponseTemplate()
		respTemplateStr := api.ResponseBody()
		headers := api.ResponseHeaders()
		status := api.Status()

		for _, flow := range ac.GetFlows() {
			var buf bytes.Buffer
			err := flow.IdentifierTemplate().Funcs(core.GetRespRenderFunctions(ac, "RO")).Execute(&buf, data)
			if err == nil {
				res := buf.String()
				if res == "true" {
					if r, ok := flow.Config()[api.Method()+" "+api.Route()]; ok {
						if r.ResponseBody != "" {
							resp, _, err = renderTemplate("body", r.ResponseBody, false)
							if err != nil {
								renderErr("body", r.ResponseBody, err)
								return
							}
							respTemplateStr = r.ResponseBody
						}
						if r.Status != "" {
							_, res, err = renderTemplate("status", r.Status, true)
							if err != nil {
								renderErr("status", r.Status, err)
								return
							}
							//parseTo int
							s, err := strconv.ParseInt(res, 10, 32)
							if err != nil {
								renderErr("status", res, err)
								return
							}
							status = int(s)
						}
						if r.Header != "" {
							_, res, err = renderTemplate("header", r.Header, true)
							if err != nil {
								renderErr("header", r.Header, err)
								return
							}
							//parseTo int
							h := common.JSONSimpleStrDict{}
							err := json.Unmarshal([]byte(res), &h)
							if err != nil {
								renderErr("header", res, err)
								return
							}
							headers = h
						}
						flow.Accessed()
						break
					}
				}
			}
		}

		//==== Response

		var buf bytes.Buffer
		err := resp.Execute(&buf, data)
		if err != nil {
			renderErr("body", respTemplateStr, err)
			return
		}
		//headers
		for k, v := range headers {
			context.Header(k, v)
		}
		context.Status(status)
		_, err = context.Writer.Write(buf.Bytes())
		if err != nil {
			ac.Log.Errorln("Error while writing to context writer", err)
		}
		return
	}
}

func routerValidator(handler gin.HandlerFunc, keyMap map[string]string) gin.HandlerFunc {
	return func(context *gin.Context) {
		for k, v := range keyMap {
			if context.Param(k) != v {
				context.AbortWithStatus(http.StatusNotFound)
				return
			}
		}
		handler(context)
	}
}

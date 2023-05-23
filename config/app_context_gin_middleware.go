package config

import (
	//"fmt"
	"mocker/common"
	"runtime/debug"

	//"mocker/error_event"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func AppContextGinMiddleware(appConfig *AppConfig) gin.HandlerFunc {
	return func(context *gin.Context) {
		requestId := common.Guid()
		ac := appConfig.NewContext(requestId)
		m := logrus.Fields{}
		xRequestId := context.GetHeader(common.XRequestIdHeaderKey)
		if xRequestId != "" {
			m[common.XRequestIdHeaderKey] = xRequestId
		}
		ac.Set("clientIp", context.ClientIP())
		params := context.Params
		for _, param := range params {
			m[param.Key] = param.Value
		}
		headers := context.Request.Header
		for k, v := range headers {
			k = "h_" + k
			m[k] = v
		}
		ac.Set("headers", headers)
		ac.Log = ac.Log.WithFields(m)
		defer func() {
			if r := recover(); r != nil {
				a := debug.Stack()
				ac.Log.Errorln("Panic recovered in request", r, string(a))
				//error_event.AddErrorEvent(ac, "panic", error_event.SEVERITY_PANIC_RECOVERED, context.Request.URL.String(), requestId, "context-middleware", fmt.Sprintln("panic:", r, string(a)), "", "")
				ac.RollbackAndLog(&ErrDbNotCommitted{}, "AppContextGinMiddleware-panic")
				panic(r)
			} else {
				ac.RollbackAndLog(&ErrDbNotCommitted{}, "AppContextGinMiddleware")
			}
		}()
		context.Set(AppContextGinContextKey, ac)
		context.Next()

	}
}

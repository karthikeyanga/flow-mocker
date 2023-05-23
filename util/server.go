package util

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mocker/common"
	"mocker/config"
	"net"
	"net/http"
	"net/url"

	"gopkg.in/yaml.v2"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var APIMODE_GIN_MODE = map[common.APIModeType]string{
	common.APIMODE_DEBUG: gin.DebugMode,
	common.APIMODE_QA:    gin.DebugMode,
	common.APIMODE_PROD:  gin.ReleaseMode,
	common.APIMODE_TEST:  gin.TestMode,
}

const (
	GIN_LOGGER_NO_LOG_KEY        = "nolog"
	GIN_LOGGER_NO_LOG_LEVEL_DATA = "data"
	GIN_LOGGER_NO_LOG_LEVEL_ALL  = "all"
	GIN_SAFE_LOG_PARAM_KEY       = "safe_log_param"
)

// BindRequestForApiResponse will try to get the info from get, post json in a greedy fashion
func BindRequestForApiResponse(context *gin.Context, requestBinder interface{}) bool {
	acceptType := context.GetHeader("Accept")
	switch acceptType {
	case "text/yaml", "text/x-yaml", "application/x-yaml", "text/yml", "text/x-yml", "application/x-yml":
		context.Header("Content-Type", acceptType)
	default:
		context.Header("Content-Type", "application/json; charset=utf-8")
	}
	if err := context.Bind(requestBinder); err != nil {
		e := &common.ErrorType{
			Code:        common.ERR_CODE_BAD_REQUEST,
			Description: "Request is missing params",
		}
		ctx := config.GetAppContext(context)
		context.Error(err).SetType(gin.ErrorTypeBind)
		if ctx.Config.APIMode == common.APIMODE_PROD {
			config.GetAppContext(context).Log.Errorln("Binding failed with error", err)
		}
		context.JSON(http.StatusBadRequest, common.BaseResponse{
			ResponseStatus: common.RespStatus_ERROR,
			Error:          e,
		})
		return false
	}
	return true
}

//QueryBindRequestForJSONResponse will populate the fields from the get params alone
func QueryBindRequestForJSONResponse(context *gin.Context, requestBinder interface{}) bool {
	context.Header("Content-Type", "application/json; charset=utf-8")
	if err := context.BindQuery(requestBinder); err != nil {
		ctx := config.GetAppContext(context)
		context.Error(err).SetType(gin.ErrorTypeBind)
		if ctx.Config.APIMode == common.APIMODE_PROD {
			config.GetAppContext(context).Log.Errorln("Binding failed with error", err)
			err = nil
		}
		context.JSON(http.StatusBadRequest, common.BaseResponse{
			ResponseStatus: common.RespStatus_ERROR,
			Error: &common.ErrorType{
				Code:        common.ERR_CODE_BAD_REQUEST,
				Description: "Request is missing params",
			},
		})
		return false
	}
	return true
}

// TODO: Should not log the card details and other sensitive information
func BindRequestForStringResponse(context *gin.Context, requestBinder interface{}) bool {
	context.Header("Content-Type", "text/plain; charset=utf-8")
	if err := context.Bind(requestBinder); err != nil {
		context.Error(err).SetType(gin.ErrorTypeBind)
		context.AbortWithStatus(http.StatusBadRequest)
		return false
	}
	return true
}

type ginRequestResponseLogWriter struct {
	gin.ResponseWriter
	responseBody *bytes.Buffer
}

func (w ginRequestResponseLogWriter) Write(b []byte) (int, error) {
	w.responseBody.Write(b)
	return w.ResponseWriter.Write(b)
}

func GinBodyLogMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		grrlw := &ginRequestResponseLogWriter{
			responseBody:   bytes.NewBufferString(""),
			ResponseWriter: context.Writer,
		}
		context.Writer = grrlw
		context.Next()
		reqBody, _ := ioutil.ReadAll(context.Request.Body)
		reqHeader, _ := json.Marshal(context.Request.Header)
		reqForm, _ := json.Marshal(context.Request.Form)
		reqPostForm, _ := json.Marshal(context.Request.PostForm)
		//statusCode := context.Writer.Status()
		response := grrlw.responseBody.String()
		logflag := true
		if nolog, present := context.Get(GIN_LOGGER_NO_LOG_KEY); present {
			switch nolog {
			case GIN_LOGGER_NO_LOG_LEVEL_ALL:
				logflag = false
			case GIN_LOGGER_NO_LOG_LEVEL_DATA:
				response = "<hidden>"
				reqForm = []byte("<hidden>")
				reqBody = reqForm
				reqPostForm = reqForm
			}
		}
		if logflag /*|| statusCode >= 400 */ {
			//ok this is an request with error, let's make a record for it
			// now print body (or log in your preferred way)
			ac := config.GetAppContext(context)
			ac.Log.WithFields(log.Fields{
				"t":         "rrl",
				"req":       string(reqBody),
				"method":    context.Request.Method,
				"resp":      response,
				"reqheader": string(reqHeader),
				"reqform":   string(reqForm),
				"reqpform":  string(reqPostForm),
				"uri":       context.Request.RequestURI,
				"addr":      context.Request.RemoteAddr,
			}).Infoln()

		}
	}
}

func GetOutboundIP() net.IP {
	//https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
	conn, err := net.Dial("udp", "8.8.8.8:80") //This doesnot connect to the address as its UDP.
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

type StringReadCloser struct {
	io.Reader
}

func (StringReadCloser) Close() error { return nil }

func GetBody(context *gin.Context) ([]byte, error) {
	body, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		return body, err
	}
	context.Request.Body = StringReadCloser{bytes.NewBuffer(body)}
	return body, nil
}

func NormalizeUrlValues(form *url.Values) map[string]interface{} {
	normalisedMap := map[string]interface{}{}
	for key, value := range *form {
		if len(value) == 1 {
			normalisedMap[key] = value[0]
		} else {
			normalisedMap[key] = value
		}
	}
	return normalisedMap
}

func SetTypedApiSuccessResponse(context *gin.Context, obj interface{}) {
	resp := common.BaseSuccessResponse{
		BaseResponse: common.BaseResponse{
			ResponseStatus: common.RespStatus_SUCCESS,
		},
		Data: obj,
	}
	setResonseByAcceptHeader(context, http.StatusOK, resp)
}

func setResonseByAcceptHeader(context *gin.Context, status int, resp interface{}) {
	acceptType := context.GetHeader("Accept")
	switch acceptType {
	case "text/yaml", "text/x-yaml", "application/x-yaml", "text/yml", "text/x-yml", "application/x-yml":
		b, err := yaml.Marshal(resp)
		if err != nil {
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		_, _ = context.Writer.Write(b)
	default:
		context.JSON(status, resp)
	}
}

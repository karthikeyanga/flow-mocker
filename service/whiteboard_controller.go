package service

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"mocker/common"
	"mocker/config"
	"mocker/core"
	"mocker/util"

	"github.com/gin-gonic/gin"
)

func WhiteboardAction(context *gin.Context) {
	ac := config.GetAppContext(context)

	defer context.Request.Body.Close()
	b, _ := ioutil.ReadAll(context.Request.Body)
	t := string(b)

	funcMap := core.GetRespRenderFunctions(ac, "RW")
	temp, err := template.New("whiteboard").Funcs(funcMap).Parse(t)
	if err != nil {
		util.SetErrorResponseFromError(context, err, "whiteboard - parsing template")
		return
	}
	var buf bytes.Buffer
	err = temp.Execute(&buf, common.ResponseVariables{})
	if err != nil {
		util.SetErrorResponseFromError(context, err, "whiteboard - parsing template")
		return
	}
	_, _ = context.Writer.Write(buf.Bytes())

}

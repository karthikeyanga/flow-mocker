package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mocker/common"
	"mocker/config"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func GetStackTrace(skip int) string {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", GetFunction(pc), GetSource(lines, line))
	}
	return string(buf.Bytes())
}

func GetFunction(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

func GetSource(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

func GetDisplayableError(ac *config.AppContext, err *common.ErrUserDisplayable) (string, string) {
	p := "general"
	genericProblem := "Something is Technically wrong - We'll fix it up real soon."
	genericSolution := "Please try again later."

	var gp, gs string
	if err != nil {
		p = err.ErrCode()
		gp = err.Problem()
		gs = err.Solution()
		if p == "" {
			return gp, gs
		}

	}
	//genericProblem, ok := ac.Config.GetString("error." + p + ".problem")
	//if !ok {
	//	if gp != "" {
	//		genericProblem = gp
	//	}
	//}
	//
	//genericSolution, ok = ac.Config.GetString("error." + p + ".solution")
	//if !ok {
	//	if gs != "" {
	//		genericSolution = gs
	//	}
	//}
	return genericProblem, genericSolution
}

func SetErrorResponseFromError(context *gin.Context, err error, where string) bool {
	ac := config.GetAppContext(context)
	if err == nil {
		return false
	}
	switch v := err.(type) {
	case *common.ErrRedirectToUrl:
		context.Redirect(http.StatusFound, v.Url)
	case *common.ErrMultiErrors:
		errs := []common.ErrorType{}
		e := common.ErrorType{
			Code:        common.ERR_CODE_VALIDATION,
			Description: "",
		}

		for f, m := range v.Errors {
			errs = append(errs, common.ErrorType{
				Code:        common.ErrorCodeType(f),
				Description: m.Error(),
			})
		}
		e.Errors = &errs

		setResonseByAcceptHeader(context, http.StatusOK, common.BaseDataResponse{
			BaseResponse: common.BaseResponse{
				ResponseStatus: common.RespStatus_ERROR,
				Error:          &e,
			},
			Data: v,
		})
	case *common.ErrUserDisplayable:
		ec := common.ERR_CODE_DISPLAYABLE
		if v.HasTobeShowedInPage() {
			ec = common.ERR_CODE_DISPLAYABLE_PAGE
		}
		problem, solution := GetDisplayableError(ac, v)
		setResonseByAcceptHeader(context, http.StatusOK, common.BaseResponse{
			ResponseStatus: common.RespStatus_ERROR,
			Error: &common.ErrorType{
				Code:        ec,
				Description: problem + " " + solution,
			},
		})
		ac.Log.Errorln("Error while", where, err)

	case *common.ErrBadRequest:
		setResonseByAcceptHeader(context, http.StatusOK, common.BaseResponse{
			ResponseStatus: common.RespStatus_ERROR,
			Error: &common.ErrorType{
				Code:        common.ERR_CODE_BAD_REQUEST,
				Description: v.Msg,
			},
		})
		ac.Log.Errorln("Error while", where, err)
	case *common.ErrUnAuthorised:
		setResonseByAcceptHeader(context, http.StatusOK, common.BaseResponse{
			ResponseStatus: common.RespStatus_ERROR,
			Error: &common.ErrorType{
				Code:        common.ERR_CODE_FORBIDDEN,
				Description: v.Msg,
			},
		})
		ac.Log.Errorln("Error while", where, err)
	default:
		ac.Log.Errorln("Error while", where, err)
		setResonseByAcceptHeader(context, http.StatusOK, common.BaseDataResponse{
			BaseResponse: common.BaseResponse{
				ResponseStatus: common.RespStatus_ERROR,
				Error: &common.ErrorType{
					Code:        common.ERR_CODE_INTERNAL_ERROR,
					Description: "Error while " + where + " : " + err.Error(),
				},
			},
			Data: err,
		})
	}
	return true
}

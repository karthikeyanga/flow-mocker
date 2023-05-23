package util

import (
	"fmt"
	"io"
	"mocker/common"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func ModAccessLogger(accessLog io.Writer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			srw := statusResponseWriter{ResponseWriter: w}
			next.ServeHTTP(&srw, r)
			rc := srw.GetStatus()
			size := srw.length
			WriteAccessLog(accessLog, r, rc, size, startTime, *srw.firstByteTime)
		})
	}
}

func AccessLogGinMiddleware(accessLog io.Writer) gin.HandlerFunc {
	return func(context *gin.Context) {
		startTime := time.Now()
		context.Next()
		size := context.Writer.Size()
		status := context.Writer.Status()
		WriteAccessLog(accessLog, context.Request, status, size, startTime, startTime)
	}
}

type statusResponseWriter struct {
	http.ResponseWriter
	status        int
	length        int
	firstByteTime *time.Time
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusResponseWriter) GetStatus() int {
	if w.firstByteTime == nil {
		now := time.Now()
		w.firstByteTime = &now
	}
	if w.status == 0 {
		return http.StatusOK
	}
	return w.status
}

func (w *statusResponseWriter) Write(b []byte) (int, error) {
	if w.length == 0 && w.firstByteTime == nil {
		now := time.Now()
		w.firstByteTime = &now
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

func GetAccessLogInfoFromRequest(r *http.Request) (clientIp, lastProxyIp, method, httpVersion, uri string) {
	clientIp = r.Header.Get("X-Forwarded-For")
	if clientIp == "" {
		clientIp = r.RemoteAddr
	}
	method = r.Method
	httpVersion = r.Proto
	uri = r.RequestURI
	lastProxyIp = r.RemoteAddr
	return
}

func WriteAccessLog(w io.Writer, r *http.Request, status, size int, startTime, firstByteTime time.Time) {
	clientIp, lastClientIp, rMethod, httpVersion, uri := GetAccessLogInfoFromRequest(r)
	tsFormat := "02/Jan/2006:15:04:05 -0700"
	now := time.Now()
	commitDuration := time.Since(firstByteTime)
	dur := time.Since(startTime)
	xRequestId := getFieldForAccessLog(r.Header.Get(common.XRequestIdHeaderKey))
	host := getFieldForAccessLog(r.Host)
	urlHost := getFieldForAccessLog(r.URL.Host)

	logLine := fmt.Sprintf("%s - - [%s] \"%s %s %s\" %d %d %d %d - %s %s %s %s\n", clientIp, now.Format(tsFormat),
		rMethod, uri, httpVersion, status, size, dur.Nanoseconds()/10e6, commitDuration.Nanoseconds()/10e6, lastClientIp, host, urlHost, xRequestId)
	if _, err := w.Write([]byte(logLine)); err != nil {
		fmt.Println("Error writing to access log : ", err, logLine)
	}
}

func getFieldForAccessLog(value string) string {
	if value == "" {
		value = "-"
	}
	return value
}

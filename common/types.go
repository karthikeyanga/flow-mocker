package common

import (
	"io"
	"net/http"
	"net/url"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type ReaderAtSeekerCloser interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
}

type AppContexter interface {
	DB() *gorm.DB
	Go(f func(AppContexter))
	Close()
	GoSync(f func(AppContexter) (interface{}, error)) (interface{}, error)
	Logger() *logrus.Entry
	Get(key string) (interface{}, bool)
	Set(key string, val interface{})
	GetRequestId() string
	RollbackAndLog(rollbackErr error, where string)
	Commit() error
	GetFlows() []IFlow
	LoadFlow(flow IFlow)
	UnloadFlow(rflow IFlow)
	EvictOldFlows()
	GetObject(id int64) (IObject, bool)
	SetObject(id int64, obj IObject)
	UnloadObject(id int64)
}

type ResponseVariables struct {
	Method     string
	Header     http.Header
	JSON       Json
	Query      url.Values
	PathParams map[string]string
	Form       url.Values
}

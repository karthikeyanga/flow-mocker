package config

import (
	"context"
	"database/sql"
	"mocker/common"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jinzhu/gorm"

	"github.com/sirupsen/logrus"
)

type ErrDbNotCommitted struct {
}

func (e *ErrDbNotCommitted) Error() string {
	return "DB not commited"
}

type AppContext struct {
	context.Context
	*AppConfig
	Log        *logrus.Entry
	curDbTrans *gorm.DB
	Now        time.Time
	requestId  string
	data       map[string]interface{}
	//parentContext *AppContext
	nestingLevel int
	threadCount  int64
	closed       bool
}

func (ctx *AppContext) Close() {
	if ctx.closed {
		return
	}
	ctx.closed = true
	//TODO: KARTHIK fix cyclic calls
	ctx.RollbackAndLog(&ErrDbNotCommitted{}, "Close")
	ctx.Log = nil
}

func (ctx *AppContext) NewContext(requestId string) *AppContext {
	nestingLevel := ctx.nestingLevel + 1
	nctx := AppContext{
		Context:   nil,
		AppConfig: ctx.AppConfig,
		Log:       ctx.Log.WithFields(map[string]interface{}{"nestedId-" + common.IntToString(nestingLevel): requestId}),
		Now:       time.Now(),
		requestId: requestId,
		data:      map[string]interface{}{},
		//parentContext: ctx,
		nestingLevel: nestingLevel,
	}
	return &nctx
}

func (ctx *AppContext) Go(f func(common.AppContexter)) {
	//We have to create a new context
	tid := atomic.AddInt64(&ctx.threadCount, 1)
	nctx := ctx.NewContext(common.Int64ToString(tid))
	for key, value := range ctx.data {
		nctx.data[key] = value
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				a := debug.Stack()
				nctx.Log.Errorln("Panic recovered in thread", r, string(a))
				//error_event.AddErrorEvent(nctx, "panic", error_event.SEVERITY_PANIC_RECOVERED, nctx.requestId, ctx.requestId, "context-Go", string(a), "", "")
			}

		}()
		f(nctx)
		nctx.RollbackAndLog(&ErrDbNotCommitted{}, "Close of Thread")
	}()
}

func (ctx *AppContext) GetRequestId() string {
	return ctx.requestId
}

func (ctx *AppContext) Set(key string, val interface{}) {
	ctx.data[key] = val
}

func (ctx *AppContext) Get(key string) (interface{}, bool) {
	v, ok := ctx.data[key]
	return v, ok
}

func (ctx *AppContext) Logger() *logrus.Entry {
	return ctx.Log
}

func (ctx *AppContext) DB() *gorm.DB {
	if ctx.closed {
		return nil
	}
	if ctx.curDbTrans == nil {
		if err := ctx.beginTransaction(); err != nil {
			ctx.Log.Errorln("Error while creating transaction", err)
		}
	}
	return ctx.curDbTrans
}

func (ctx *AppContext) beginTransaction() error {
	if ctx.curDbTrans == nil {
		ctx.Log.Warningln("BEGIN DB TRANSACTION...")
		db := ctx.db.Debug().BeginTx(context.Background(), &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		})
		ctx.curDbTrans = db
		if err := db.Error; err != nil {
			return err
		}
	}
	return nil
}
func (ctx *AppContext) RollbackAndLog(rollbackErr error, where string) {
	close := false
	if ctx.curDbTrans != nil {
		if _, ok := rollbackErr.(*ErrDbNotCommitted); ok {
			ctx.Log.Warningln("Rolling back from AppContextGinMiddleware/Thread : Please check if this is intended and you did not forget to commit")
			close = true
		}
		if err := ctx.curDbTrans.Debug().Rollback().Error; err != nil {
			ctx.Log.WithFields(logrus.Fields{
				"function": where,
			}).Errorln("Error while rollback", rollbackErr)
			ctx.curDbTrans = nil
			return
		}
		ctx.curDbTrans = nil
		if close {
			ctx.Close()
		}
	}
}
func (ctx *AppContext) Commit() error {
	if ctx.curDbTrans != nil {
		ctx.Log.Infoln("COMMITTING... ")
		if err := ctx.curDbTrans.Debug().Commit().Error; err != nil {
			if err == sql.ErrTxDone {
				ctx.Log.Warningln("Transaction Done error", err)
			} else {
				ctx.curDbTrans = nil
				return err
			}
		}
		ctx.curDbTrans = nil
	}
	return nil
}

func GetAppContext(ginCtx *gin.Context) *AppContext {
	ctx, ok := ginCtx.Get(AppContextGinContextKey)
	if ok {
		return ctx.(*AppContext)
	}
	return nil
}

/**
Example:
result,err := ctx.GoSync(func(ac *AppContext) (interface{}, error) {
	return map[string]string{}, nil
})
*/
func (ctx *AppContext) GoSync(f func(common.AppContexter) (interface{}, error)) (interface{}, error) {
	//We have to create a new context
	tid := atomic.AddInt64(&ctx.threadCount, 1)
	nctx := ctx.NewContext("Sync-" + common.Int64ToString(tid))
	defer nctx.Close()
	for key, value := range ctx.data {
		nctx.data[key] = value
	}
	return f(nctx)
}

func (ctx *AppContext) HandlePanic(src string, onPanic func(interface{}), finally func()) {
	if r := recover(); r != nil {
		a := debug.Stack()
		ctx.Log.Errorln("Panic recovered in default handler", r, string(a))
		//error_event.AddErrorEvent(ctx, "panic", error_event.SEVERITY_PANIC_RECOVERED, ctx.requestId,
		//	ctx.requestId, "context-Go", string(a), "", "")
		if onPanic != nil {
			onPanic(r)
		}
	}
	if finally != nil {
		finally()
	}
}

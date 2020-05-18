package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"runtime/debug"
	"sync"

	"github.com/CharlesBases/common/log"
)

const (
	nilRequestMessage = "Request is nil"
)

const (
	UserErrorStopRun = "user stop run"
	UserErrorMySQL   = "mysql connect error"
	UserErrorRedis   = "redis connect error"
	UserErrorNSQ     = "nsq connect error"
)

var (
	usererror = map[string]interface{}{
		UserErrorStopRun: nil,
		UserErrorMySQL:   nil,
		UserErrorRedis:   nil,
		UserErrorNSQ:     nil,
	}
)

type panicRecover struct {
	PrintStack       bool
	StackAll         bool
	StackSize        int
	PanicHandlerFunc func(*panicInformation)
}

func Recovery() *panicRecover {
	return &panicRecover{
		PrintStack: true,
		StackAll:   false,
		StackSize:  1024 * 8,
	}
}

func (rec *panicRecover) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {

			swg := sync.WaitGroup{}
			swg.Add(1)
			defer swg.Wait()

			go func() {
				defer swg.Done()
				systemError(rw)
			}()

			if err, ok := usererror[fmt.Sprintf("%v", err)]; ok {
				log.Errorf("Panic: [%s]", err)
				return
			}

			stack := make([]byte, rec.StackSize)
			stack = stack[:runtime.Stack(stack, rec.StackAll)]
			infor := &panicInformation{RecoveredPanic: err, Request: r}

			if rec.PrintStack {
				infor.Stack = stack
			}
			log.Errorf("Panic: [%s]\nRequestInfo: [%s]\nStackInfo:   %s", infor.RecoveredPanic, infor.RequestDescription(), infor.StackAsString())

			if rec.PanicHandlerFunc != nil {
				func() {
					defer func() {
						if err := recover(); err != nil {
							log.Error(fmt.Sprintf("provided PanicHandlerFunc panic'd: %s, trace:\n%s", err, debug.Stack()))
							log.Error(fmt.Sprintf("%s\n", debug.Stack()))
						}
					}()
					rec.PanicHandlerFunc(infor)
				}()
			}
		}
	}()

	next(rw, r)
}

type panicInformation struct {
	RecoveredPanic interface{}
	Stack          []byte
	Request        *http.Request
}

func (p *panicInformation) StackAsString() string {
	return string(p.Stack)
}

func (p *panicInformation) RequestDescription() string {
	if p.Request == nil {
		return nilRequestMessage
	}

	var queryOutput string
	if len(p.Request.URL.RawQuery) != 0 {
		queryOutput = "?" + p.Request.URL.RawQuery
	}
	return fmt.Sprintf("%s %s%s", p.Request.Method, p.Request.URL.Path, queryOutput)
}

func systemError(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusInternalServerError)
	data, _ := json.Marshal(map[string]interface{}{
		"errNo":  500,
		"errMsg": "请求错误",
	})
	rw.Write(data)
}

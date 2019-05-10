package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"runtime/debug"

	"github.com/CharlesBases/common/log"
)

const (
	nilRequestMessage = "Request is nil"
)

var (
	errAbort = errors.New("user stop run")
)

type PanicRecover struct {
	PrintStack       bool
	StackAll         bool
	StackSize        int
	PanicHandlerFunc func(*PanicInformation)
}

func (rec *PanicRecover) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			if err == errAbort {
				// 用户主动结束
				return
			}

			stack := make([]byte, rec.StackSize)
			stack = stack[:runtime.Stack(stack, rec.StackAll)]
			infor := &PanicInformation{RecoveredPanic: err, Request: r}

			if rec.PrintStack {
				infor.Stack = stack
			}

			log.Errorf("Panic - RequestInfo = %s ; Err = %s ; StackInfo = %s", infor.RequestDescription(), infor.RecoveredPanic, infor.StackAsString())

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

			rw.WriteHeader(http.StatusInternalServerError)
			msg, _ := json.Marshal(&struct {
				code int
				msg  string
			}{
				521,
				"服务器错误",
			})
			rw.Write(msg)
		}
	}()
	next(rw, r)
}

func Recovery() *PanicRecover {
	return &PanicRecover{
		PrintStack: true,
		StackAll:   false,
		StackSize:  1024 * 8,
	}
}

// PanicInformation contains all
// elements for printing stack informations.
type PanicInformation struct {
	RecoveredPanic interface{}
	Stack          []byte
	Request        *http.Request
}

// StackAsString returns a printable version of the stack
func (p *PanicInformation) StackAsString() string {
	return string(p.Stack)
}

// RequestDescription returns a printable description of the url
func (p *PanicInformation) RequestDescription() string {
	if p.Request == nil {
		return nilRequestMessage
	}

	var queryOutput string
	if p.Request.URL.RawQuery != "" {
		queryOutput = "?" + p.Request.URL.RawQuery
	}
	return fmt.Sprintf("%s %s%s", p.Request.Method, p.Request.URL.Path, queryOutput)
}

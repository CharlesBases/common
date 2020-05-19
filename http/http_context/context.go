package http_context

import (
	"time"
)

type httpContext struct {
	*logger

	Request  *[]byte
	Response interface{}

	Communication chan struct{}
}

func NewContext() *httpContext {
	return &httpContext{logger: newlogger()}
}

func (httpContext *httpContext) ThrowCheck(err error, errmsg string) {
	if err != nil {

	}
}

func (httpContext *httpContext) Monitor() {
}

func (httpContext *httpContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (httpContext *httpContext) Value(key interface{}) interface{} {
	return nil
}

func (httpContext *httpContext) Done() <-chan struct{} {
	return nil
}

func (httpContext *httpContext) Err() error {
	return nil
}

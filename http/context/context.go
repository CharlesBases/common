package context

import (
	"context"
	"time"
)

type Context struct {
	context.Context

	request  *[]byte
	response interface{}

	Infors map[interface{}]interface{}

	logger *logger
}

func NewContext() *Context {
	return &Context{logger: newlogger()}
}

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (ctx *Context) Done() <-chan struct{} {
	return nil
}

func (ctx *Context) Err() error {
	return nil
}

func (ctx *Context) Value(key interface{}) interface{} {
	return nil
}

func (ctx *Context) Flush() {
	ctx.logger.flush()
}

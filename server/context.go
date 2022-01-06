package server

import (
	"github.com/go-xxl/xxl/log"
	"github.com/go-xxl/xxl/utils"
	"net/http"
	"sync"
	"time"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Param   *RunReq
	TraceId string
	m       sync.Map
}

// Context impl

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
	val, ok := ctx.m.Load(key)
	if ok {
		return val
	}
	return nil
}

func (ctx *Context) String() string {
	return ""
}

func (ctx *Context) SetTraceId(traceId string) {
	ctx.TraceId = traceId
}

func (ctx *Context) GetTraceId() string {
	return ctx.TraceId
}

func (ctx *Context) Copy() *Context {
	c := &Context{
		Writer:  ctx.Writer,
		Request: ctx.Request,
		Param:   ctx.Param,
		TraceId: ctx.TraceId,
		m:       sync.Map{},
	}

	ctx.m.Range(func(key, value interface{}) bool {
		c.m.Store(key, value)
		return true
	})

	return c
}

// resp

func (ctx *Context) JSON(code int, data []byte) {
	_, _ = ctx.Writer.Write(data)
}

func (ctx *Context) Success(msg string, data interface{}) {
	resp := Resp{
		Code:    SuccessCode,
		Msg:     msg,
		Content: data,
	}
	log.Info(msg, log.Field("data", utils.ObjToStr(data)))
	ctx.JSON(SuccessCode, utils.ObjToBytes(resp))
}

func (ctx *Context) Fail(msg string, data interface{}) {
	resp := Resp{
		Code:    FailCode,
		Msg:     msg,
		Content: data,
	}
	log.Info(msg, log.Field("data", utils.ObjToStr(data)))
	ctx.JSON(FailCode, utils.ObjToBytes(resp))
}

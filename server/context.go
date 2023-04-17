package server

import (
	"context"
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
	return ctx.Request.Context().Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.Request.Context().Done()
}

func (ctx *Context) Err() error {
	return ctx.Request.Context().Err()
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
    r, e := http.NewRequestWithContext(context.Background(), ctx.Request.Method, ctx.Request.RequestURI, ctx.Request.Body)
    if e != nil {
       r = ctx.Request
    }
    r.Clone(ctx.Request.Context())

	c := &Context{
		Writer:  ctx.Writer,
		Request: r,
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
	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Writer.WriteHeader(code)
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

///////////////// sync map set

func (ctx *Context) Set(key string, value interface{}) {
	ctx.m.Store(key, value)
}

func (ctx *Context) Get(key string) (value interface{}, exist bool) {
	return ctx.m.Load(key)
}

func (ctx *Context) GetString(key string) string {
	if val, ok := ctx.Get(key); ok {
		if value, isString := val.(string); isString {
			return value
		}
	}
	return ""
}

func (ctx *Context) GetStringOrDefault(key, defaultValue string) string {
	if val, ok := ctx.Get(key); ok {
		if value, isString := val.(string); isString {
			return value
		}
	}
	return defaultValue
}

func (ctx *Context) GetInt64(key string) int64 {
	if val, ok := ctx.Get(key); ok {
		if value, isInt64 := val.(int64); isInt64 {
			return value
		}
	}
	return 0
}

func (ctx *Context) GetInt64OrDefault(key string, defaultValue int64) int64 {
	if val, ok := ctx.Get(key); ok {
		if value, isInt64 := val.(int64); isInt64 {
			return value
		}
	}
	return defaultValue
}

// Reset reset ctx
func (ctx *Context) Reset() *Context {
	ctx.Writer = nil
	ctx.Request = nil
	ctx.Param = nil
	ctx.TraceId = ""
	ctx.m = sync.Map{}
	return ctx
}

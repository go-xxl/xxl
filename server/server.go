package server

import (
	"bytes"
	"encoding/json"
	"github.com/go-xxl/xxl/log"
	"github.com/go-xxl/xxl/utils"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type Handler func(ctx *Context)

type Engine struct {
	pool        sync.Pool
	mu          sync.Mutex
	funcHandler map[string]Handler
	before      []Handler
	after       []Handler
}

func New() *Engine {
	engine := &Engine{
		pool:        sync.Pool{},
		funcHandler: make(map[string]Handler),
	}
	engine.pool.New = func() interface{} {
		return engine.allocateContext()
	}

	return engine
}

func (engine *Engine) allocateContext() *Context {
	return &Context{}
}

func (engine *Engine) Run(addr ...string) error {
	address := utils.ResolveAddress(addr)
	localSrv := utils.GetPort(address)
	log.Info("Starting server at " + localSrv)
	return http.ListenAndServe(":"+localSrv, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := engine.pool.Get().(*Context)
	c.Writer = w
	c.Request = req
	c.TraceId = utils.Uuid()

	engine.handleHTTPRequest(c)

	engine.pool.Put(c)
}

func (engine *Engine) handleHTTPRequest(ctx *Context) {

	b, _ := ioutil.ReadAll(ctx.Request.Body)
	_ = json.Unmarshal(b, &ctx.Param)

	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	rPath := ctx.Request.URL
	if handler, ok := engine.funcHandler[rPath.Path]; ok {
		engine.beforeHandler(ctx)
		handler(ctx)
		engine.afterHandler(ctx)
		return
	}

	http.NotFound(ctx.Writer, ctx.Request)
	return
}

// AddRoute add web endpoint
func (engine *Engine) AddRoute(path string, handler Handler) {
	engine.funcHandler[engine.formatRoute(path)] = handler
}

func (engine *Engine) formatRoute(path string) string {
	if strings.HasPrefix(path, "/") {
		return path
	}
	return "/" + path
}

func (engine *Engine) beforeHandler(ctx *Context) {
	if len(engine.before) > 0 {
		for _, handler := range engine.before {
			handler(ctx)
		}
	}
}

func (engine *Engine) SetBeforeHandlers(handlers ...Handler) {
	engine.before = append(engine.before, handlers...)
}

func (engine *Engine) afterHandler(ctx *Context) {
	if len(engine.after) > 0 {
		for _, handler := range engine.after {
			handler(ctx)
		}
	}
}

func (engine *Engine) SetAfterHandler(handlers ...Handler) {
	engine.after = append(engine.after, handlers...)
}

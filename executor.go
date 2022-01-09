package xxl

import (
	"github.com/go-xxl/xxl/admin"
	"github.com/go-xxl/xxl/job"
	"github.com/go-xxl/xxl/server"
	"github.com/go-xxl/xxl/utils"
	"log"
	"net/http"
	"net/http/pprof"
)

type JobExecutor interface {
	// Job add a job endpoint
	Job(handlerName string, handler job.Func)

	WithHealthCheck(path string, handler server.Handler)

	WithDebug(isDebug bool)

	Run()
}

type HealthFunc func(ctx *server.Context)

// Executor Executor
type Executor struct {
	opts    Options
	address string
	engine  *server.Engine
	isPprof bool
}

// NewExecutor create a JobExecutor
func NewExecutor(opts ...Option) JobExecutor {
	opt := Options{
		ExecutorIp:   utils.GetLocalIp(),
		ExecutorPort: DefaultExecutorPort,
		RegistryKey:  DefaultRegistryKey,
	}

	for _, o := range opts {
		o(&opt)
	}

	e := &Executor{
		opts: opt,
	}

	e.address = utils.BuildEndPoint(opt.ExecutorIp, opt.ExecutorPort)
	e.opts.RegistryValue = e.address
	e.engine = server.New()
	return e
}

// Run start job service
func (e *Executor) Run() {

	e.engine.SetBeforeHandlers(func(ctx *server.Context) {

	})

	e.engine.SetAfterHandler(func(ctx *server.Context) {

	})

	biz := NewExecutorService()
	e.engine.AddRoute("/run", biz.Run)
	e.engine.AddRoute("/kill", biz.Kill)
	e.engine.AddRoute("/log", biz.Log)
	e.engine.AddRoute("/beat", biz.Beat)
	e.engine.AddRoute("/idleBeat", biz.IdleBeat)

	if e.isPprof {
		e.engine.AddRoute("/debug/pprof/", e.WrapF(pprof.Index))
		e.engine.AddRoute("/debug/pprof/cmdline", e.WrapF(pprof.Cmdline))
		e.engine.AddRoute("/debug/pprof/profile", e.WrapF(pprof.Profile))
		e.engine.AddRoute("/debug/pprof/symbol", e.WrapF(pprof.Symbol))
		e.engine.AddRoute("/debug/pprof/trace", e.WrapF(pprof.Trace))

		e.engine.AddRoute("/debug/pprof/heap", e.WrapF(pprof.Index))
		e.engine.AddRoute("/debug/pprof/goroutine", e.WrapF(pprof.Index))
		e.engine.AddRoute("/debug/pprof/allocs", e.WrapF(pprof.Index))
		e.engine.AddRoute("/debug/pprof/block", e.WrapF(pprof.Index))
		e.engine.AddRoute("/debug/pprof/mutex", e.WrapF(pprof.Index))
		e.engine.AddRoute("/debug/pprof/threadcreate", e.WrapF(pprof.Index))
	}

	adm := admin.NewAdmApi()
	adm.SetOpt(
		admin.SetAccessToken(e.opts.AccessToken),
		admin.SetAdmAddresses(e.opts.AdmAddresses),
		admin.SetTimeout(e.opts.Timeout),
		admin.SetRegistryKey(e.opts.RegistryKey),
		admin.SetRegistryValue(e.opts.RegistryValue),
		admin.SetRegistryGroup(admin.GetGroupName(admin.EXECUTOR)),
	)

	defer adm.RegistryRemove()

	go func() {
		adm.Register()
		log.Fatalln(e.engine.Run(e.address))
	}()

	utils.WatchSignal()
}

// WithHealthCheck ExecutorService's web health check endpoint
func (e *Executor) WithHealthCheck(path string, handler server.Handler) {
	if path != "" {
		e.engine.AddRoute(path, handler)
	}
}

// WithDebug open pprof
func (e *Executor) WithDebug(debug bool) {
	e.isPprof = debug
}

func (e *Executor) Job(handlerName string, handler job.Func) {
	job.GetAllHandlerList().Set(handlerName, job.Job{
		Id:        0,
		Name:      "",
		Ext:       nil,
		Cancel:    nil,
		Param:     nil,
		Fn:        handler,
		StartTime: 0,
		EndTime:   0,
	})
}

func (e *Executor) WrapF(f func(w http.ResponseWriter, r *http.Request)) server.Handler {
	return func(ctx *server.Context) {
		f(ctx.Writer, ctx.Request)
	}
}

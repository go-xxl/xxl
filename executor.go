package xxl

import (
	"github.com/go-xxl/xxl/admin"
	"github.com/go-xxl/xxl/job"
	"github.com/go-xxl/xxl/server"
	"github.com/go-xxl/xxl/utils"
	"log"
)

type JobExecutor interface {
	// Job add a job endpoint
	Job(handlerName string, handler job.Func)

	WithHealthCheck(path string, handler server.Handler)

	Run()
}

type HealthFunc func(ctx *server.Context)

// Executor Executor
type Executor struct {
	opts    Options
	address string
	engine  *server.Engine
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

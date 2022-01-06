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

	Run()
}

// Executor Executor
type Executor struct {
	ctx     server.Context
	opts    Options
	address string
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
	return e
}

// Run start job service
func (e *Executor) Run() {
	engine := server.New()
	engine.SetBeforeHandlers(func(ctx *server.Context) {})

	engine.SetAfterHandler(func(ctx *server.Context) {})

	biz := NewExecutorService()
	engine.AddRoute("/run", biz.Run)
	engine.AddRoute("/kill", biz.Kill)
	engine.AddRoute("/log", biz.Log)
	engine.AddRoute("/beat", biz.Beat)
	engine.AddRoute("/idleBeat", biz.IdleBeat)

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
		log.Fatalln(engine.Run(e.address))
	}()

	utils.WatchSignal()
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

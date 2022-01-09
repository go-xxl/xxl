package main

import (
	"github.com/go-xxl/xxl"
	"github.com/go-xxl/xxl/admin"
	"github.com/go-xxl/xxl/job"
	"github.com/go-xxl/xxl/server"
	"github.com/go-xxl/xxl/utils"
	"time"
)

func main() {

	//log.SetLog(&CLog{})

	e := xxl.NewExecutor(
		xxl.AdmAdmAddresses("http://192.168.0.104:8080/xxl-job-admin/"),
		xxl.ExecutorIp(utils.GetLocalIp()),
		xxl.ExecutorPort("12345"),
		xxl.RegistryKey("demo-test"),
	)

	e.Job("/demo", func(ctx *server.Context) job.Resp {
		param := ctx.Param

		time.Sleep(time.Second * 30)

		return job.Resp{
			LogID:       param.LogID,
			LogDateTime: time.Now().Unix(),
			HandleCode:  admin.SuccessCode,
			HandleMsg:   "get result",
		}
	})

	e.WithHealthCheck("/health", func(ctx *server.Context) {
		ctx.Success("pong return", "pong")
		return
	})

	e.WithDebug(true)

	e.Run()
}

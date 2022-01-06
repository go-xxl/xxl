# xxl-job client for Golang

# Installation

```go
go get github.com/go-xxl/xxl
```

# Quickstart

```
package main

import (
	"github.com/go-xxl/xxl"
	"github.com/go-xxl/xxl/admin"
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

	e.Job("/demo", func(ctx *server.Context) server.JobResp {
		param := ctx.Param

		time.Sleep(time.Second * 30)

		return server.JobResp{
			LogID:       param.LogID,
			LogDateTime: time.Now().Unix(),
			HandleCode:  admin.SuccessCode,
			HandleMsg:   "get result",
		}
	})

	e.Run()
}
```

# Tree
```go
.
├── _example
├── admin
├── job
├── log
├── server
├── utils
└── vendor
    └── go.uber.org
        ├── atomic
        ├── multierr
        └── zap
            ├── buffer
            ├── internal
            │   ├── bufferpool
            │   ├── color
            │   └── exit
            └── zapcore
```
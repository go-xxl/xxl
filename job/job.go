package job

import (
	"context"
	"github.com/go-xxl/xxl/server"
	"sync"
)

type Func func(ctx *server.Context) Resp

type (
	Job struct {
		Id        int64
		Name      string
		Ext       context.Context
		Cancel    context.CancelFunc
		Fn        Func
		Param     *server.RunReq
		StartTime int64
		EndTime   int64
	}

	Resp struct {
		LogID       int64  `json:"logId"`
		LogDateTime int64  `json:"logDateTime"`
		HandleCode  int    `json:"handleCode"`
		HandleMsg   string `json:"handleMsg"`
	}
)

var (
	RunningHandlerList *HandlerList
	AllHandlerList     *HandlerList
	runningOnce        sync.Once
	allOnce            sync.Once
)

func GetRunningJobList() *HandlerList {
	runningOnce.Do(func() {
		RunningHandlerList = &HandlerList{}
	})
	return RunningHandlerList
}

func GetAllHandlerList() *HandlerList {
	allOnce.Do(func() {
		AllHandlerList = &HandlerList{}
	})
	return AllHandlerList
}

type HandlerList struct {
	sync.Map
}

func (t *HandlerList) Set(key string, val Job) {
	t.Store(key, val)
}

func (t *HandlerList) Get(key string) (Job, bool) {
	job, ok := t.Load(key)
	if ok {
		return job.(Job), ok
	}
	return Job{}, false
}

func (t *HandlerList) Del(key string) {
	t.Delete(key)
}

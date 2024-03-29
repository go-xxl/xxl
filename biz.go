package xxl

import (
	"context"
	"encoding/json"
	"github.com/go-xxl/xxl/admin"
	"github.com/go-xxl/xxl/job"
	"github.com/go-xxl/xxl/log"
	"github.com/go-xxl/xxl/server"
	"github.com/go-xxl/xxl/utils"
	"go.uber.org/zap"
	"io/ioutil"
	"time"
)

type ExecutorBiz struct {
}

func NewExecutorService() *ExecutorBiz {
	return &ExecutorBiz{}
}

// Run start web service
func (biz *ExecutorBiz) Run(ctx *server.Context) {
	param := ctx.Param
	var (
		jobHandler job.Job
		exist      bool
	)
	if jobHandler, exist = job.GetAllHandlerList().Get(param.ExecutorHandler); !exist {
		ctx.Fail("job is not exist, job's id is "+utils.Int2Str(param.JobID)+", job's name is "+param.ExecutorHandler, "")
		return
	}

	var (
		runningJob   job.Job
		runningExist bool
	)
	if runningJob, runningExist = job.GetRunningJobList().Get(utils.Int2Str(param.JobID)); runningExist {
		if param.ExecutorBlockStrategy != CoverEarly {
			ctx.Fail("job is already running, job's id is ["+utils.Int2Str(param.JobID)+"], job's name is "+param.ExecutorHandler, "")
			return
		}

		runningJob.Cancel()
		job.GetRunningJobList().Del(utils.Int2Str(param.JobID))
	}

	taskCtx := ctx.Copy()
	if param.ExecutorTimeout > 0 {
		jobHandler.Ext, jobHandler.Cancel = context.WithTimeout(taskCtx.Request.Context(), time.Duration(param.ExecutorTimeout)*time.Second)
	} else {
		jobHandler.Ext, jobHandler.Cancel = context.WithCancel(taskCtx.Request.Context())
	}


    taskCtx.Request = taskCtx.Request.WithContext(jobHandler.Ext)

	jobHandler.Id = param.JobID
	jobHandler.Name = param.ExecutorHandler
	jobHandler.Param = param
	jobHandler.StartTime = time.Now().Unix()

	job.GetRunningJobList().Set(utils.Int2Str(jobHandler.Id), jobHandler)

	go func() {
		defer func() {
			if e := recover(); e != nil {
				log.GetLog().Error("task panic", zap.Any("e", e))
			}
		}()

		defer func() {
			job.GetRunningJobList().Del(utils.Int2Str(jobHandler.Id))
		}()

		var resp job.Resp
		defer func() {
			var call []admin.HandleCallbackParams
			call = append(call, admin.HandleCallbackParams{
				LogId:      resp.LogID,
				LogDateTim: resp.LogDateTime,
				HandleCode: resp.HandleCode,
				HandleMsg:  resp.HandleMsg,
			})
			admin.GetClient().Callback(call)
		}()

		resp = jobHandler.Fn(taskCtx)
	}()

	traceId := ctx.TraceId
	ctx.Success("job is starting …… , job's id is "+utils.Int2Str(param.JobID)+", job's name is "+param.ExecutorHandler+"，job's traceId is "+traceId, "")
}

// Kill stop job
func (biz *ExecutorBiz) Kill(ctx *server.Context) {
	param := ctx.Param

	runningJob, exist := job.GetRunningJobList().Get(utils.Int2Str(param.JobID))

	if !exist {
		ctx.Fail("job is not running ……, job's id is "+utils.Int2Str(runningJob.Id)+", job's name is "+runningJob.Name, "")
		return
	}
	runningJob.Cancel()
	job.GetRunningJobList().Del(utils.Int2Str(param.JobID))
	ctx.Success("job is removed, job's id is "+utils.Int2Str(runningJob.Id)+", job's name is "+runningJob.Name, "")
}

// Log job log
func (biz *ExecutorBiz) Log(ctx *server.Context) {
	var resp LogResp
	data, _ := ioutil.ReadAll(ctx.Request.Body)
	req := &LogReq{}
	if err := json.Unmarshal(data, &req); err != nil {
		ctx.Fail("params err", resp)
		return
	}
	resp.LogContent = "The distributed system does not store logs, please go to the log service to check."
	resp.IsEnd = false
	resp.FromLineNum = req.FromLineNum
	resp.ToLineNum = req.FromLineNum + 2
	ctx.Success("Log parsing request parsing completed", resp)
	return
}

// Beat check alive
func (biz *ExecutorBiz) Beat(ctx *server.Context) {
	ctx.Success("beat ing...", nil)
}

// IdleBeat check job is alive
func (biz *ExecutorBiz) IdleBeat(ctx *server.Context) {
	param := ctx.Param

	_, exist := job.GetRunningJobList().Get(utils.Int2Str(param.JobID))

	if exist {
		ctx.Fail("idleBeat is running, job's id is "+utils.Int2Str(param.JobID), "")
		return
	}
	ctx.Success("idle beat ing...", param)
}

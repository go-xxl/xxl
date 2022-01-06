package server

type RunReq struct {
	JobID                 int64  `json:"jobId"`
	ExecutorHandler       string `json:"executorHandler"`
	ExecutorParams        string `json:"executorParams"`
	ExecutorBlockStrategy string `json:"executorBlockStrategy"`
	ExecutorTimeout       int64  `json:"executorTimeout"`
	LogID                 int64  `json:"logId"`
	LogDateTime           int64  `json:"logDateTime"`
	GlueType              string `json:"glueType"`
	GlueSource            string `json:"glueSource"`
	GlueUpdatetime        int64  `json:"glueUpdatetime"`
	BroadcastIndex        int64  `json:"broadcastIndex"`
	BroadcastTotal        int64  `json:"broadcastTotal"`
}

type Resp struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Content interface{} `json:"content"`
}

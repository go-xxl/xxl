package xxl


type LogResp struct {
	FromLineNum int    `json:"fromLineNum"`
	ToLineNum   int    `json:"toLineNum"`
	LogContent  string `json:"logContent"`
	IsEnd       bool   `json:"isEnd"`
}

type LogReq struct {
	LogDateTim  int64 `json:"logDateTim"`
	LogId       int   `json:"logId"`
	FromLineNum int   `json:"fromLineNum"`
}

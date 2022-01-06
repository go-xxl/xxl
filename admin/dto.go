package admin

type Resp struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Content interface{} `json:"content"`
}

type Registry struct {
	RegistryGroup string `json:"registryGroup"`
	RegistryKey   string `json:"registryKey"`
	RegistryValue string `json:"registryValue"`
}

type HandleCallbackParams struct {
	LogId      int64  `json:"logId"`
	LogDateTim int64  `json:"logDateTim"`
	HandleCode int    `json:"handleCode"`
	HandleMsg  string `json:"handleMsg"`
}

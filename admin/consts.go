package admin

import "time"

const (
	SuccessCode = 200
	FailCode    = 500
)

const (
	EmptyString = ""

	UrlRegister       = "api/registry"
	UrlRegistryRemove = "api/registryRemove"
	UrlCallback       = "api/callback"

	XxlJobAccessToken = "XXL-JOB-ACCESS-TOKEN"
)

// RegistryConfig
const (
	EXECUTOR = iota
	ADMIN
)

func GetGroupName(groupId int) string {
	switch groupId {
	case EXECUTOR:
		return "EXECUTOR"
	case ADMIN:
		return "ADMIN"
	default:
		return ""
	}
}

const LoopFrequency = time.Second * 30

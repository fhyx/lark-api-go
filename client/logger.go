package client

import (
	"fhyx.online/lark-api-go/log"
)

func logger() log.Logger {
	return log.GetLogger()
}

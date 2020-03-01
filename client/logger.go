package client

import (
	"github.com/fhyx/lark-api-go/log"
)

func logger() log.Logger {
	return log.GetLogger()
}

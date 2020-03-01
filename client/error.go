package client

import (
	"fmt"
)

// Error ...
type Error struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"msg,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

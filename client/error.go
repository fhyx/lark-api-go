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
	return fmt.Sprintf("code: %d, msg: %s", e.Code, e.Message)
}

// GetCode ...
func (e *Error) GetCode() int {
	return e.Code
}

// GetMessage ...
func (e *Error) GetMessage() string {
	return e.Message
}

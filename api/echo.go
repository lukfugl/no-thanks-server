package api

import (
	"fmt"
)

// An EchoAction is an example ApiAction that simply logs the provided message
type EchoAction struct {
	Message string `json:"message"`
}

// Execute logs the provided message
func (a *EchoAction) Execute() error {
	fmt.Printf("echo: %s\n", a.Message)
	return nil
}

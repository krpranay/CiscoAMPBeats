// This package will be used to log different types of error and information

package logmanager

import (
	"fmt"
)

// Error type
type Error error

// Enum created for different error type
const (
	INFORMATION = "INFORMATION"
	WARNING     = "WARNING"
	ERROR       = "ERROR"
)

// Log function to log error
func Log(errorType string, customError string, exMessage Error) {
	fmt.Println(errorType+" : ", customError, exMessage.Error())

}

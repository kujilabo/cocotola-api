package domain

import (
	"errors"
	"strings"
)

var ErrPluginError = errors.New("plugin error")

type PluginError struct {
	ErrorType     string
	ErrorCode     string
	ErrorMessages []string
	OrgError      error
}

func NewPluginError(errorType, errorCode string, errorMessages []string, err error) *PluginError {
	return &PluginError{
		ErrorType:     errorType,
		ErrorCode:     errorCode,
		ErrorMessages: errorMessages,
		OrgError:      err,
	}
}

func (e *PluginError) Error() string {
	return strings.Join(e.ErrorMessages, "\n")
}

// type PluginClientError interface {
// 	Error() string
// 	GetErrorCode() string
// 	GetErrorMessages() []string
// }

// type pluginClientError struct {
// 	ErrorCode     string
// 	ErrorMessages []string
// }

// func NewPluginClientError(errorCode string, errorMessages []string) PluginClientError {
// 	return &pluginClientError{
// 		ErrorCode:     errorCode,
// 		ErrorMessages: errorMessages,
// 	}
// }

// func (e *pluginClientError) GetErrorCode() string {
// 	return e.ErrorCode
// }

// func (e *pluginClientError) GetErrorMessages() []string {
// 	return e.ErrorMessages
// }

// func (e *pluginClientError) Error() string {
// 	return strings.Join(e.ErrorMessages, "\n")
// }

// type PluginServerError interface {
// 	Error() string
// 	GetErrorCode() string
// 	GetErrorMessages() []string
// }

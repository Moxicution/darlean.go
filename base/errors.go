package base

import (
	"fmt"
	"strings"

	"github.com/darlean-io/darlean.go/utils/variant"
)

type ErrorKind string

const ERROR_KIND_FRAMEWORK = "framework"
const ERROR_KIND_APPLICATION = "application"

type ActionError struct {
	Code       string         `json:"code"`
	Message    string         `json:"message"`
	Template   string         `json:"template"`
	Kind       ErrorKind      `json:"kind"`
	Parameters map[string]any `json:"parameters"`
	Nested     []*ActionError `json:"nested"`
	Stack      string         `json:"stack"`
}

type ActionErrorOptions struct {
	Code       string
	Template   string
	Parameters map[string]any
	Nested     []*ActionError
	Stack      string
}

func (error ActionError) Error() string {
	return error.Message
}

func FormatTemplate(template string, parameters map[string]any) string {
	message := template
	for key, value := range parameters {
		assignable, supported := value.(variant.Assignable)
		if supported {
			var target any
			err := assignable.AssignTo(target)
			if err != nil {
				value = err.Error()
			} else {
				value = target
			}
		}
		message = strings.ReplaceAll(message, "["+key+"]", "\""+fmt.Sprint(value)+"\"")
	}
	return message
}

func newActionError(options ActionErrorOptions, kind ErrorKind) *ActionError {
	e := ActionError{
		Kind:       kind,
		Code:       options.Code,
		Message:    options.Template,
		Template:   options.Template,
		Parameters: options.Parameters,
		Nested:     options.Nested,
		Stack:      options.Stack,
	}
	if len(e.Parameters) > 0 {
		e.Message = FormatTemplate(options.Template, e.Parameters)
	}
	if e.Code != "" {
		if e.Message == "" {
			e.Message = e.Code
		} else {
			e.Message = "(" + e.Code + ") " + e.Message
		}
	}
	return &e
}

func NewFrameworkError(options ActionErrorOptions) *ActionError {
	return newActionError(options, ERROR_KIND_FRAMEWORK)
}

func NewApplicationError(options ActionErrorOptions) *ActionError {
	return newActionError(options, ERROR_KIND_APPLICATION)
}
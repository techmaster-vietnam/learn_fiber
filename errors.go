package main

import (
	"fmt"
)

// ============================================================================
// Custom Error Types - FRAMEWORK AGNOSTIC
// ============================================================================
type ErrorType string

const (
	BusinessError   ErrorType = "BUSINESS"   // 4xx
	SystemError     ErrorType = "SYSTEM"     // 5xx
	ValidationError ErrorType = "VALIDATION" // 400
	AuthError       ErrorType = "AUTH"       // 401-403
	ExternalError   ErrorType = "EXTERNAL"   // 502-504
	PanicError      ErrorType = "PANIC"      // recovered panic
)

type AppError struct {
	Type      ErrorType
	Code      int
	Message   string
	Details   map[string]interface{}
	Cause     error
	RequestID string
}

func (e *AppError) Error() string {
	return e.Message
}

// ============================================================================
// Factory Functions - Tạo Error Dễ Dàng (FRAMEWORK AGNOSTIC)
// ============================================================================

func NewBusinessError(code int, msg string) *AppError {
	file, line, function := getCallerInfo(1)
	return &AppError{
		Type:    BusinessError,
		Code:    code,
		Message: msg,
		Details: map[string]interface{}{
			"function": function,
			"file":     fmt.Sprintf("%s:%d", file, line),
		},
	}
}

func NewSystemError(err error) *AppError {
	file, line, function := getCallerInfo(1)
	return &AppError{
		Type:    SystemError,
		Code:    500,
		Message: "Internal server error",
		Cause:   err,
		Details: map[string]interface{}{
			"function": function,
			"file":     fmt.Sprintf("%s:%d", file, line),
		},
	}
}

func NewValidationError(msg string, details map[string]interface{}) *AppError {
	file, line, function := getCallerInfo(1)
	if details == nil {
		details = make(map[string]interface{})
	}
	details["function"] = function
	details["file"] = fmt.Sprintf("%s:%d", file, line)
	return &AppError{Type: ValidationError, Code: 400, Message: msg, Details: details}
}

func NewAuthError(code int, msg string) *AppError {
	file, line, function := getCallerInfo(1)
	return &AppError{
		Type:    AuthError,
		Code:    code,
		Message: msg,
		Details: map[string]interface{}{
			"function": function,
			"file":     fmt.Sprintf("%s:%d", file, line),
		},
	}
}

func NewExternalError(code int, msg string, cause error) *AppError {
	file, line, function := getCallerInfo(1)
	return &AppError{
		Type:    ExternalError,
		Code:    code,
		Message: msg,
		Cause:   cause,
		Details: map[string]interface{}{
			"function": function,
			"file":     fmt.Sprintf("%s:%d", file, line),
		},
	}
}

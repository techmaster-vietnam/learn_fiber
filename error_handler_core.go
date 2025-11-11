package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// ============================================================================
// Core Error Handler Logic - FRAMEWORK AGNOSTIC
// ============================================================================

// HandlePanic xử lý panic và trả về AppError
func HandlePanic(r interface{}, requestID string) *AppError {
	actualFile, actualLine, actualFunc := getActualPanicLocation()
	callChain := formatStackTraceArray()

	return &AppError{
		Type:      PanicError,
		Code:      500,
		Message:   fmt.Sprintf("Panic recovered: %v", r),
		RequestID: requestID,
		Details: map[string]interface{}{
			"panic_value": r,
			"function":    actualFunc,
			"file":        fmt.Sprintf("%s:%d", actualFile, actualLine),
			"call_chain":  callChain,
		},
	}
}

// ConvertToAppError chuyển đổi error thường thành AppError
func ConvertToAppError(err error, requestID string) *AppError {
	// Check nếu đã là AppError
	if appErr, ok := err.(*AppError); ok {
		appErr.RequestID = requestID
		return appErr
	}

	// Convert error thường thành AppError
	return &AppError{
		Type:      SystemError,
		Code:      500,
		Message:   "Internal server error",
		Cause:     err,
		RequestID: requestID,
	}
}

// ============================================================================
// Log và Response Helper - FRAMEWORK AGNOSTIC
// ============================================================================

// LogError xử lý logging cho error
// Sử dụng dual-logger strategy:
// - Console: log tất cả (development)
// - File: chỉ log lỗi nghiêm trọng (Panic, System, External)
func LogError(appErr *AppError, requestPath string) {
	// Chuẩn bị log fields
	fields := logrus.Fields{
		"error_type":  string(appErr.Type),
		"status_code": appErr.Code,
		"path":        requestPath,
		"request_id":  appErr.RequestID,
	}

	// Thêm details nếu có
	for k, v := range appErr.Details {
		fields[k] = v
	}

	// Thêm cause nếu có
	if appErr.Cause != nil {
		fields["cause"] = appErr.Cause.Error()
	}

	// Log vào FILE (chỉ lỗi nghiêm trọng)
	if isSevereError(appErr.Type) {
		fileLogger.WithFields(fields).Error(appErr.Message)
	}
}

// FormatErrorResponse tạo response data cho client
func FormatErrorResponse(appErr *AppError) map[string]interface{} {
	return map[string]interface{}{
		"error":      appErr.Message,
		"type":       string(appErr.Type),
		"request_id": appErr.RequestID,
	}
}

// LogAndRespond xử lý logging và gửi response (framework agnostic)
func LogAndRespond(ctx HTTPContext, appErr *AppError, requestPath string) {
	// 1. Log error
	LogError(appErr, requestPath)

	// 2. Send response
	ctx.Status(appErr.Code).JSON(FormatErrorResponse(appErr))
}

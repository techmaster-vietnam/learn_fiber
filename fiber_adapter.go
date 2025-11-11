package main

import (
	"github.com/gofiber/fiber/v2"
)

// ============================================================================
// Fiber Adapter - Implement HTTPContext Interface
// ============================================================================

// FiberContext wrap Fiber's context để implement HTTPContext interface
type FiberContext struct {
	ctx *fiber.Ctx
}

// NewFiberContext tạo FiberContext từ fiber.Ctx
func NewFiberContext(c *fiber.Ctx) *FiberContext {
	return &FiberContext{ctx: c}
}

func (f *FiberContext) Method() string {
	return f.ctx.Method()
}

func (f *FiberContext) Path() string {
	return f.ctx.Path()
}

func (f *FiberContext) GetLocal(key string) interface{} {
	return f.ctx.Locals(key)
}

func (f *FiberContext) Status(code int) HTTPContext {
	f.ctx.Status(code)
	return f
}

func (f *FiberContext) JSON(data interface{}) error {
	return f.ctx.JSON(data)
}

// ============================================================================
// Fiber Error Handler Middleware
// ============================================================================

// FiberErrorHandlerMiddleware là Fiber-specific middleware
// Nó wrap core error handling logic
func FiberErrorHandlerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Wrap Fiber context
		ctx := NewFiberContext(c)

		requestPath := ctx.Method() + " " + ctx.Path()
		requestID := "unknown"
		if rid, ok := ctx.GetLocal("requestid").(string); ok {
			requestID = rid
		}

		// Panic recovery
		defer func() {
			r := recover()
			if r != nil {
				// Xử lý panic bằng core logic
				panicErr := HandlePanic(r, requestID)
				LogAndRespond(ctx, panicErr, requestPath)
			}
		}()

		// Thực thi handler
		err := c.Next()

		// Xử lý error nếu có
		if err != nil {
			// Convert sang AppError bằng core logic
			appErr := ConvertToAppError(err, requestID)
			LogAndRespond(ctx, appErr, requestPath)
			return nil
		}

		return nil
	}
}

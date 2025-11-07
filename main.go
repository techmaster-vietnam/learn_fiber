package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
)

// ============================================================================
// Global Logger Instance
// ============================================================================
var (
	logrusLogger *logrus.Logger
)

// init khởi tạo logrus logger
func init() {
	initLogrus()
}

// initLogrus khởi tạo logrus logger
func initLogrus() {
	logrusLogger = logrus.New()
	logrusLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		PrettyPrint:     true,
	})
	logrusLogger.SetReportCaller(false) // Tắt caller information tự động để tránh trùng lặp
	logrusLogger.SetLevel(logrus.DebugLevel)
}

func main() {
	app := fiber.New(fiber.Config{
		AppName: "LearnFiber - Logrus Demo",
	})

	// Middleware
	app.Use(logger.New())

	// Routes
	app.Get("/", homeHandler)
	app.Get("/logrus", logrusHandler)
	app.Get("/logrus2", logrus2Handler)
	app.Get("/logrus3", logrus3Handler)

	// Start server
	fmt.Println("Server starting on http://localhost:8081")
	if err := app.Listen(":8081"); err != nil {
		panic(err)
	}
}

func homeHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Welcome to LearnFiber - Logrus Demo",
		"endpoints": []string{
			"/logrus - Demo sirupsen/logrus (chia cho 0)",
			"/logrus2 - Demo sirupsen/logrus (index out of range)",
			"/logrus3 - Demo sirupsen/logrus (deep call stack: X->Y->Z->GetElement)",
		},
	})
}

// ============================================================================
// sirupsen/logrus Handler
// ============================================================================
func logrusHandler(c *fiber.Ctx) error {
	requestPath := fmt.Sprintf("%s %s", c.Method(), c.Path())

	// Defer function để bắt panic từ chia cho 0
	defer func() {
		if r := recover(); r != nil {
			// Lấy thông tin CHÍNH XÁC về dòng gây panic
			file, line, funcName := getPanicLocationForHandler("logrusHandler")

			// Log panic với thông tin tối giản
			logrusLogger.WithFields(logrus.Fields{
				"request_path": requestPath,
				"function":     funcName,
				"file":         fmt.Sprintf("%s:%d", file, line),
			}).Error("PANIC: Lỗi chia cho 0 đã xảy ra!")
		}
	}()

	// Dòng này sẽ gây panic!
	denominator := 0
	result := 100 / denominator

	return c.JSON(fiber.Map{"result": result})
}

// ============================================================================
// sirupsen/logrus Handler 2 - Index Out of Range
// ============================================================================
func logrus2Handler(c *fiber.Ctx) error {
	requestPath := fmt.Sprintf("%s %s", c.Method(), c.Path())

	// Defer function để bắt panic từ index out of range
	defer func() {
		if r := recover(); r != nil {
			// Lấy thông tin về dòng THỰC SỰ gây panic (tận cùng)
			actualFile, actualLine, actualFunc := getActualPanicLocation()

			// Format stack trace thành dạng dễ đọc
			callChainArray := formatStackTraceArray()

			// Log panic với thông tin chi tiết
			logrusLogger.WithFields(logrus.Fields{
				"panic_value":  fmt.Sprintf("%v", r),
				"function":     actualFunc,
				"location":     fmt.Sprintf("%s:%d", actualFile, actualLine),
				"request_path": requestPath,
				"call_chain":   callChainArray,
			}).Error("PANIC: Lỗi truy cập mảng ngoài phạm vi!")
		}
	}()

	// Dòng này sẽ gây panic!
	element := GetElement()

	return c.JSON(fiber.Map{"element": element})
}

// GetElement truy cập phần tử mảng không tồn tại
func GetElement() int {
	arr := []int{1, 2, 3}
	return arr[10] // Index out of range - panic tại đây!
}

// ============================================================================
// sirupsen/logrus Handler 3 - Deep Call Stack Test
// ============================================================================
func logrus3Handler(c *fiber.Ctx) error {
	requestPath := fmt.Sprintf("%s %s", c.Method(), c.Path())

	// Defer function để bắt panic từ deep call stack
	defer func() {
		if r := recover(); r != nil {
			// Lấy thông tin về dòng THỰC SỰ gây panic (tận cùng)
			actualFile, actualLine, actualFunc := getActualPanicLocation()

			// Format stack trace thành dạng dễ đọc
			callChainArray := formatStackTraceArray()

			logrusLogger.WithFields(logrus.Fields{
				"panic_value":  fmt.Sprintf("%v", r),
				"function":     actualFunc,
				"location":     fmt.Sprintf("%s:%d", actualFile, actualLine),
				"request_path": requestPath,
				"call_chain":   callChainArray,
			}).Error("PANIC: Test deep call stack - X -> Y -> Z -> GetElement")
		}
	}()

	// Chuỗi gọi hàm nhiều tầng
	result := callX()

	return c.JSON(fiber.Map{"result": result})
}

// callX gọi callY
func callX() int {
	return callY()
}

// callY gọi callZ
func callY() int {
	return callZ()
}

// callZ gọi GetElement
func callZ() int {
	return callW() // Panic sẽ xảy ra trong GetElement, không phải ở đây!
}

// callZ gọi GetElement
func callW() int {
	return GetElement() // Panic sẽ xảy ra trong GetElement, không phải ở đây!
}

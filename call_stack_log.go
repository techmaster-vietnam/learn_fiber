package main

import (
	"fmt"
	"runtime/debug"
	"strings"
)

// getActualPanicLocation lấy thông tin về dòng THỰC SỰ gây panic (frame đầu tiên sau panic)
// Đây là nơi thực sự phát sinh lỗi, không phải nơi gọi hàm
func getActualPanicLocation() (file string, line int, function string) {
	// Lấy stack trace từ debug.Stack()
	stack := string(debug.Stack())
	lines := strings.Split(stack, "\n")

	// Tìm dòng "panic" trong stack trace
	// Stack trace format:
	// ...
	// panic({...})
	//     /path/to/panic.go:XXX
	// main.GetElement(...)         <- Frame này là nơi thực sự gây panic!
	//     /path/to/main.go:216
	// main.logrus2Handler(...)
	//     /path/to/main.go:208

	panicFound := false
	for i := 0; i < len(lines); i++ {
		l := strings.TrimSpace(lines[i])

		// Tìm dòng "panic"
		if !panicFound && strings.HasPrefix(l, "panic(") {
			panicFound = true
			continue
		}

		// Sau khi tìm thấy panic, bỏ qua dòng location của panic function
		// Dòng tiếp theo là function thực sự gây panic
		if panicFound && strings.HasPrefix(l, "main.") {
			// Đây là function thực sự gây panic
			function = l
			// Lấy tên function (bỏ phần parameter)
			if idx := strings.Index(function, "("); idx > 0 {
				function = function[:idx]
			}

			// Dòng tiếp theo chứa file:line
			if i+1 < len(lines) {
				locationLine := strings.TrimSpace(lines[i+1])

				// Parse file và line
				// Format: /path/to/main.go:XXX +0x...
				parts := strings.Fields(locationLine)
				if len(parts) > 0 {
					fileAndLine := parts[0]
					if idx := strings.LastIndex(fileAndLine, ":"); idx > 0 {
						file = fileAndLine[:idx]
						fmt.Sscanf(fileAndLine[idx+1:], "%d", &line)
					}
				}
			}
			break
		}
	}

	if file == "" {
		return "unknown", 0, "unknown"
	}

	return file, line, function
}

// getPanicLocationForHandler lấy thông tin chính xác về dòng gây panic từ stack trace
// handlerName: tên handler như "logrusHandler", "slogHandler", "zapHandler", "zerologHandler"
func getPanicLocationForHandler(handlerName string) (file string, line int, function string) {
	// Lấy stack trace từ debug.Stack()
	stack := string(debug.Stack())
	lines := strings.Split(stack, "\n")

	// Tìm pattern "main.<handlerName>(" nhưng KHÔNG phải "main.<handlerName>.func1"
	pattern := "main." + handlerName + "("

	// Tìm từ đầu stack trace
	for i := 0; i < len(lines); i++ {
		l := strings.TrimSpace(lines[i])

		// Tìm pattern nhưng không phải func1 (defer function)
		if strings.HasPrefix(l, pattern) && !strings.Contains(l, "func1") {
			function = "main." + handlerName

			// Dòng tiếp theo chứa file:line
			if i+1 < len(lines) {
				locationLine := strings.TrimSpace(lines[i+1])

				// Parse file và line
				// Format: /path/to/main.go:XXX +0x...
				parts := strings.Fields(locationLine)
				if len(parts) > 0 {
					fileAndLine := parts[0]
					if idx := strings.LastIndex(fileAndLine, ":"); idx > 0 {
						file = fileAndLine[:idx]
						fmt.Sscanf(fileAndLine[idx+1:], "%d", &line)
					}
				}
			}
			break
		}
	}

	if file == "" {
		return "unknown", 0, "unknown"
	}

	return file, line, function
}

// formatStackTraceArray format stack trace thành array dễ đọc
// Chỉ lấy các hàm trong package main và vị trí gọi
func formatStackTraceArray() []string {
	stack := string(debug.Stack())
	lines := strings.Split(stack, "\n")

	var callChain []string
	skipNext := false

	// Danh sách các hàm cần bỏ qua (utility functions)
	skipFunctions := map[string]bool{
		"main.formatStackTrace":           true,
		"main.formatStackTraceArray":      true,
		"main.getActualPanicLocation":     true,
		"main.getPanicLocationForHandler": true,
	}

	for i := 0; i < len(lines); i++ {
		l := strings.TrimSpace(lines[i])

		// Bỏ qua các dòng runtime internal
		if skipNext {
			skipNext = false
			continue
		}

		// Chỉ lấy các hàm trong package main
		if strings.HasPrefix(l, "main.") {
			funcName := l
			// Lấy tên function (bỏ phần parameter)
			if idx := strings.Index(funcName, "("); idx > 0 {
				funcName = funcName[:idx]
			}

			// Bỏ qua anonymous functions (có chứa .func) và utility functions
			if strings.Contains(funcName, ".func") || skipFunctions[funcName] {
				skipNext = true
				continue
			}

			// Dòng tiếp theo chứa file:line
			if i+1 < len(lines) {
				locationLine := strings.TrimSpace(lines[i+1])
				parts := strings.Fields(locationLine)
				if len(parts) > 0 {
					fileAndLine := parts[0]
					// Chỉ lấy tên file, bỏ đường dẫn đầy đủ
					if idx := strings.LastIndex(fileAndLine, "/"); idx >= 0 {
						fileAndLine = fileAndLine[idx+1:]
					}

					callChain = append(callChain, fmt.Sprintf("%s (%s)", funcName, fileAndLine))
				}
			}
			skipNext = true
		}
	}

	// Thứ tự tự nhiên của stack trace: từ nơi gây lỗi lên handler
	return callChain
}

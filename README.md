# LearnFiber - Logging Libraries Demo

Dự án demo sử dụng Fiber framework để thử nghiệm các thư viện logging phổ biến trong Go.

## Mô tả

Ứng dụng web này cho phép bạn thử nghiệm và so sánh 4 thư viện logging khác nhau:

1. **log/slog** - Thư viện logging chuẩn của Go (từ Go 1.21)
2. **sirupsen/logrus** - Structured logger phổ biến
3. **uber-go/zap** - High-performance logger từ Uber
4. **rs/zerolog** - Zero-allocation JSON logger

## Tính năng

Mỗi thư viện logging hiển thị đầy đủ thông tin:

- ✅ **Cấp độ lỗi** (Error level)
- ✅ **Thông điệp lỗi** (Error message)
- ✅ **Thông tin source code**: package, function, file, line number
- ✅ **Stack trace** đầy đủ
- ✅ **Các biến cục bộ**: user_id, user_name, request_path, request_method, etc.
- ✅ **Timestamp** với định dạng chuẩn

## Cài đặt

### Yêu cầu

- Go 1.21 trở lên
- npm (theo [[memory:8760137]])

### Các bước cài đặt

1. Clone hoặc tạo dự án:

```bash
cd /Users/cuong/CODE/LearnFiber
```

2. Cài đặt dependencies (đã được cài sẵn):

```bash
go mod download
```

3. Build ứng dụng:

```bash
go build -o learnfiber main.go
```

## Sử dụng

### Chạy server

```bash
go run main.go
```

Hoặc chạy file đã build:

```bash
./learnfiber
```

Server sẽ khởi động tại: **http://localhost:8081**

### Các Endpoints

| Endpoint | Thư viện | Mô tả |
|----------|----------|-------|
| `/` | - | Trang chủ với danh sách endpoints |
| `/slog` | log/slog | Demo logging với slog |
| `/logrus` | sirupsen/logrus | Demo logging với logrus |
| `/zap` | uber-go/zap | Demo logging với zap |
| `/zerolog` | rs/zerolog | Demo logging với zerolog |

### Ví dụ

1. Mở trình duyệt hoặc dùng curl:

```bash
# Xem danh sách endpoints
curl http://localhost:8081/

# Test log/slog
curl http://localhost:8081/slog

# Test sirupsen/logrus
curl http://localhost:8081/logrus

# Test uber-go/zap
curl http://localhost:8081/zap

# Test rs/zerolog
curl http://localhost:8081/zerolog
```

2. Kiểm tra console/terminal để xem log output chi tiết

## So sánh các thư viện

### log/slog
- ✅ Built-in, không cần dependency ngoài
- ✅ JSON handler với structured logging
- ✅ Hỗ trợ context
- ✅ AddSource option cho caller info

### sirupsen/logrus
- ✅ Structured logging với Fields
- ✅ Pretty print JSON
- ✅ Nhiều formatter có sẵn
- ✅ SetReportCaller() cho caller info

### uber-go/zap
- ✅ High-performance, zero-allocation
- ✅ Structured logging
- ✅ Tự động thêm stack trace ở error level
- ✅ Color-coded output trong development mode

### rs/zerolog
- ✅ Zero-allocation JSON logger
- ✅ Fluent API (chainable methods)
- ✅ Console writer với pretty format
- ✅ Rất nhanh và hiệu quả về memory

## Cấu trúc dự án

```
LearnFiber/
├── main.go          # File chính chứa tất cả code
├── go.mod           # Go module definition
├── go.sum           # Go dependencies checksums
├── learnfiber       # Compiled binary
└── README.md        # File này
```

## Dependencies

```
github.com/gofiber/fiber/v2 v2.52.9
github.com/sirupsen/logrus v1.9.3
go.uber.org/zap v1.27.0
github.com/rs/zerolog v1.34.0
```

## Ghi chú

- Mỗi endpoint sẽ log ra console với format riêng của từng thư viện
- Tất cả đều hiển thị error message, source location, local variables và stack trace
- Code được tổ chức rõ ràng với comment để dễ học tập và tham khảo

## License

MIT License - Dự án học tập và demo


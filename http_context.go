package main

// ============================================================================
// HTTP Context Interface - FRAMEWORK AGNOSTIC
// ============================================================================

// HTTPContext là interface chung cho mọi web framework
// Mỗi framework (Fiber, Gin, Echo, Chi...) sẽ implement interface này
type HTTPContext interface {
	// Method trả về HTTP method (GET, POST, PUT...)
	Method() string

	// Path trả về request path
	Path() string

	// GetLocal lấy giá trị từ context locals
	GetLocal(key string) interface{}

	// Status set HTTP status code
	Status(code int) HTTPContext

	// JSON trả về JSON response
	JSON(data interface{}) error
}

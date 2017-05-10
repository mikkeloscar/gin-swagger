package api

// Response is a simple response from an HTTP service.
type Response struct {
	Code int
	Body interface{}
}

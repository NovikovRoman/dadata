package client

import "fmt"

// ResponseError indicates an HTTP non-200 response code.
type ResponseError struct {
	Status     string // e.g. "200 OK"
	StatusCode int    // e.g. 200
	Body       []byte // response body (truncated), usually contains DaData error details
}

func (e *ResponseError) Error() string {
	if len(e.Body) > 0 {
		return fmt.Sprintf("HTTP response: %s: %s", e.Status, e.Body)
	}
	return fmt.Sprintf("HTTP response: %s", e.Status)
}

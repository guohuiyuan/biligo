package bilibili

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrMissingSessData = errors.New("bilibili: missing SESSDATA")
	ErrMissingBiliJct  = errors.New("bilibili: missing bili_jct")
)

// APIError represents a logical error returned by Bilibili.
type APIError struct {
	Code    int
	Message string
	Data    []byte
}

func (e *APIError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("bilibili api error: code=%d", e.Code)
	}
	return fmt.Sprintf("bilibili api error: code=%d message=%s", e.Code, e.Message)
}

// HTTPError represents a transport-level failure.
type HTTPError struct {
	StatusCode int
	Method     string
	URL        string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("bilibili http error: %s %s status=%d", e.Method, e.URL, e.StatusCode)
}

func (e *HTTPError) Temporary() bool {
	return e.StatusCode == http.StatusTooManyRequests || e.StatusCode >= 500
}

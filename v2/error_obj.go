package twitter

import "fmt"

// ResponseDecodeError is an error when a response has a decoding error, JSON.
type ResponseDecodeError struct {
	Name      string
	Err       error
	RateLimit *RateLimit
}

func (r *ResponseDecodeError) Error() string {
	return fmt.Sprintf("%s decode error: %v", r.Name, r.Err)
}

// Unwrap will return the wrapped error
func (r *ResponseDecodeError) Unwrap() error {
	return r.Err
}

// HTTPError is a response error where the body is not JSON, but XML.  This commonly seen in 404 errors.
type HTTPError struct {
	Status     string
	StatusCode int
	URL        string
	RateLimit  *RateLimit
}

func (h HTTPError) Error() string {
	return fmt.Sprintf("twitter [%s] status: %s code: %d", h.URL, h.Status, h.StatusCode)
}

// ErrorObj is part of the partial errors in the response
type ErrorObj struct {
	Title        string      `json:"title"`
	Detail       string      `json:"detail"`
	Type         string      `json:"type"`
	ResourceType string      `json:"resource_type"`
	Parameter    string      `json:"parameter"`
	Value        interface{} `json:"value"`
}

// Error is part of the HTTP response error
type Error struct {
	Parameters interface{} `json:"parameters"`
	Message    string      `json:"message"`
}

// ErrorResponse is returned by a non-success callout
type ErrorResponse struct {
	StatusCode int
	Errors     []Error    `json:"errors"`
	Title      string     `json:"title"`
	Detail     string     `json:"detail"`
	Type       string     `json:"type"`
	RateLimit  *RateLimit `json:"-"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("twitter callout status %d %s:%s", e.StatusCode, e.Title, e.Detail)
}

package twitter

import "fmt"

// HTTPError is a response error where the body is not JSON, but XML.  This commonly seen in 404 errors.
type HTTPError struct {
	Status     string
	StatusCode int
	URL        string
}

func (h *HTTPError) Error() string {
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
	Errors     []Error `json:"errors"`
	Title      string  `json:"title"`
	Detail     string  `json:"detail"`
	Type       string  `json:"type"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("twiiter callout status %d %s:%s", e.StatusCode, e.Title, e.Detail)
}

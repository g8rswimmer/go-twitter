package twitter

import "fmt"

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

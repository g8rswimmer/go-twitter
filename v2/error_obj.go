package twitter

// ErrorObj is part of the partial errors in the response
type ErrorObj struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Type   string `json:"type"`
}

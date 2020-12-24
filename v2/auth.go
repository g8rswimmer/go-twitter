package twitter

import "net/http"

// Authorizer will add the authorization to the HTTP request
type Authorizer interface {
	Add(req *http.Request)
}

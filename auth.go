package twitter

import "net/http"

type Authorizer interface {
	Add(*http.Request)
}

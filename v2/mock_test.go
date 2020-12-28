package twitter

import "net/http"

type roundTripFunc func(request *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func mockHTTPClient(fn roundTripFunc) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(fn),
	}
}

type mockAuth struct{}

func (m *mockAuth) Add(*http.Request) {}

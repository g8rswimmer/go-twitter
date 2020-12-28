package twitter

import "fmt"

type endpoint string

const (
	tweetLookupEndpoint endpoint = "2/tweets"
)

func (e endpoint) url(host string) string {
	return fmt.Sprintf("%s/%s", host, string(e))
}

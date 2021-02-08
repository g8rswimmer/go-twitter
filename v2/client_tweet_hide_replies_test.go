package twitter

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestClient_TweetHideReplies(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		id   string
		hide bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodPut {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodPut)
					}
					if strings.Contains(req.URL.String(), tweetHideRepliesEndpoint.urlID("", "63046977")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetHideRepliesEndpoint.urlID("", "63046977"))
					}
					body := `{"data":{"hidden":true}}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				id:   "63046977",
				hide: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			if err := c.TweetHideReplies(context.Background(), tt.args.id, tt.args.hide); (err != nil) != tt.wantErr {
				t.Errorf("Client.TweetHideReplies() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package twitter

import (
	"context"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestClient_QuoteTweetsLookup(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		tweetID string
		opts    QuoteTweetsLookupOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *QuoteTweetsLookupResponse
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), quoteTweetLookupEndpoint.urlID("", "tweet-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), quoteTweetLookupEndpoint)
					}
					body := `{
						"data": [
							{
								"id": "1503982413004914689",
								"text": "RT @suhemparack: Super excited to share our course on Getting started with the #TwitterAPI v2 for academic research\n\nIf you know students w…"
							}
						],
						"meta": {
							"result_count": 1,
							"next_token": "axdnchiqasch"
						}
					}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(body)),
						Header: func() http.Header {
							h := http.Header{}
							h.Add(rateLimit, "15")
							h.Add(rateRemaining, "12")
							h.Add(rateReset, "1644461060")
							return h
						}(),
					}
				}),
			},
			args: args{
				tweetID: "tweet-1234",
			},
			want: &QuoteTweetsLookupResponse{
				Raw: &TweetRaw{
					Tweets: []*TweetObj{
						{
							ID:   "1503982413004914689",
							Text: "RT @suhemparack: Super excited to share our course on Getting started with the #TwitterAPI v2 for academic research\n\nIf you know students w…",
						},
					},
				},
				Meta: &QuoteTweetsLookupMeta{
					ResultCount: 1,
					NextToken:   "axdnchiqasch",
				},
				RateLimit: &RateLimit{
					Limit:     15,
					Remaining: 12,
					Reset:     Epoch(1644461060),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := c.QuoteTweetsLookup(context.Background(), tt.args.tweetID, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.QuoteTweetsLookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.QuoteTweetsLookup() = %v, want %v", got, tt.want)
			}
		})
	}
}

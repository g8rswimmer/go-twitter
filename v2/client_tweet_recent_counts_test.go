package twitter

import (
	"context"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestClient_TweetRecentCounts(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		query string
		opts  TweetRecentCountsOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TweetRecentCountsResponse
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
					if strings.Contains(req.URL.String(), string(tweetRecentCountsEndpoint)) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetRecentCountsEndpoint)
					}
					body := `{
					  "data": [
						{
						  "end": "2021-05-27t00:00:00.000z",
						  "start": "2021-05-26t23:00:00.000z",
						  "tweet_count": 2
						},
						{
						  "end": "2021-05-27t01:00:00.000z",
						  "start": "2021-05-27t00:00:00.000z",
						  "tweet_count": 2
						}
					  ],
					  "meta": {
						"total_tweet_count": 4
					  }
					}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				query: "python",
			},
			want: &TweetRecentCountsResponse{
				TweetCounts: []*TweetCount{
					{
						End:        "2021-05-27t00:00:00.000z",
						Start:      "2021-05-26t23:00:00.000z",
						TweetCount: 2,
					},
					{
						End:        "2021-05-27t01:00:00.000z",
						Start:      "2021-05-27t00:00:00.000z",
						TweetCount: 2,
					},
				},
				Meta: &TweetRecentCountsMeta{
					TotalTweetCount: 4,
				},
			},
			wantErr: false,
		},
		{
			name: "success-optional",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), string(tweetRecentCountsEndpoint)) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetRecentCountsEndpoint)
					}
					body := `{
					  "data": [
						{
						   "start": "2021-10-08T15:29:42.000Z",
						   "end": "2021-10-09T00:00:00.000Z",
						   "tweet_count": 2
						},
					    {
						   "start": "2021-10-09T00:00:00.000Z",
						   "end": "2021-10-09T15:29:33.000Z",
						   "tweet_count": 2
					    }
					  ],
					  "meta": {
						"total_tweet_count": 4
					  }
					}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				query: "python",
				opts: TweetRecentCountsOpts{
					StartTime:   time.Now().Add(-24 * time.Hour),
					Granularity: Granularity("day"),
				},
			},
			want: &TweetRecentCountsResponse{
				TweetCounts: []*TweetCount{
					{
						End:        "2021-10-09T00:00:00.000Z",
						Start:      "2021-10-08T15:29:42.000Z",
						TweetCount: 2,
					},
					{
						End:        "2021-10-09T15:29:33.000Z",
						Start:      "2021-10-09T00:00:00.000Z",
						TweetCount: 2,
					},
				},
				Meta: &TweetRecentCountsMeta{
					TotalTweetCount: 4,
				},
			},
			wantErr: false,
		},
		{
			name: "Bad Request",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), string(tweetRecentCountsEndpoint)) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetRecentCountsEndpoint)
					}
					body := `{
						"errors": [
							{
								"parameters": {
									"id": [
										"aassd"
									]
								},
								"message": "The id query parameter value [aassd] does not match ^[0-9]{1,19}$"
							}
						],
						"title": "Invalid Request",
						"detail": "One or more parameters to your request was invalid.",
						"type": "https://api.twitter.com/2/problems/invalid-request"
					}`
					return &http.Response{
						StatusCode: http.StatusBadRequest,
						Body:       io.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				query: "nothing",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := c.TweetRecentCounts(context.Background(), tt.args.query, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.TweetRecentCounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.TweetRecentCounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

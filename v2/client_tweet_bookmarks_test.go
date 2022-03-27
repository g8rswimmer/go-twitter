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

func TestClient_TweetBookmarksLookup(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		userID string
		opts   TweetBookmarksLookupOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TweetBookmarksLookupResponse
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
					if strings.Contains(req.URL.String(), tweetBookmarksEndpoint.urlID("", "user-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetBookmarksEndpoint)
					}
					body := `{
						"data": [
						  {
							"id": "1294346980072624128",
							"text": "I awake from five years of slumber https://t.co/OEPVyAFcfB"
						  }
						],
						"meta": {
						  "result_count": 1,
						  "next_token": "zldjwdz3w6sba13nbs0mbravfipbtqvbiqplg9h0p4k"
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
				userID: "user-1234",
			},
			want: &TweetBookmarksLookupResponse{
				Raw: &TweetRaw{
					Tweets: []*TweetObj{
						{
							ID:   "1294346980072624128",
							Text: "I awake from five years of slumber https://t.co/OEPVyAFcfB",
						},
					},
				},
				Meta: &TweetBookmarksLookupMeta{
					ResultCount: 1,
					NextToken:   "zldjwdz3w6sba13nbs0mbravfipbtqvbiqplg9h0p4k",
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
			got, err := c.TweetBookmarksLookup(context.Background(), tt.args.userID, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.TweetBookmarksLookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.TweetBookmarksLookup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_AddTweetBookmark(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		userID  string
		tweetID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *AddTweetBookmarkResponse
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodPost {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodPost)
					}
					if strings.Contains(req.URL.String(), tweetBookmarksEndpoint.urlID("", "user-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetBookmarksEndpoint)
					}
					body := `{
						"data": {
						  "bookmarked": true
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
				userID:  "user-1234",
				tweetID: "tweet-5678",
			},
			want: &AddTweetBookmarkResponse{
				Tweet: &TweetBookmarkData{
					Bookmarked: true,
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
			got, err := c.AddTweetBookmark(context.Background(), tt.args.userID, tt.args.tweetID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.AddTweetBookmark() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.AddTweetBookmark() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_RemoveTweetBookmark(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		userID  string
		tweetID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *RemoveTweetBookmarkResponse
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodDelete {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodDelete)
					}
					if strings.Contains(req.URL.String(), tweetBookmarksEndpoint.urlID("", "user-1234")+"/tweet-5678") == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetBookmarksEndpoint)
					}
					body := `{
						"data": {
						  "bookmarked": false
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
				userID:  "user-1234",
				tweetID: "tweet-5678",
			},
			want: &RemoveTweetBookmarkResponse{
				Tweet: &TweetBookmarkData{
					Bookmarked: false,
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
			got, err := c.RemoveTweetBookmark(context.Background(), tt.args.userID, tt.args.tweetID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.RemoveTweetBookmark() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.RemoveTweetBookmark() = %v, want %v", got, tt.want)
			}
		})
	}
}

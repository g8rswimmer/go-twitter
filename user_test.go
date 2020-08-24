package twitter

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestUser_Lookup(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		ids       []string
		fieldOpts UserFieldOptions
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           UserLookups
		wantErr        bool
		wantTweetError *TweetErrorResponse
	}{
		{
			name: "success id",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), userLookupEndpoint) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetLookupEndpoint)
					}
					body := `{
						"data": {
						  "id": "2244994945",
						  "name": "Twitter Dev",
						  "username": "TwitterDev"
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				ids:       []string{"2244994945"},
				fieldOpts: UserFieldOptions{},
			},
			want: UserLookups{
				"2244994945": UserLookup{
					User: UserObj{
						UserName: "TwitterDev",
						ID:       "2244994945",
						Name:     "Twitter Dev",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "success ids",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), userLookupEndpoint) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetLookupEndpoint)
					}
					body := `{
						"data": [
						  {
							"created_at": "2013-12-14T04:35:55.000Z",
							"username": "TwitterDev",
							"pinned_tweet_id": "1255542774432063488",
							"id": "2244994945",
							"name": "Twitter Dev"
						  },
						  {
							"created_at": "2007-02-20T14:35:54.000Z",
							"username": "Twitter",
							"pinned_tweet_id": "1274087687469715457",
							"id": "783214",
							"name": "Twitter"
						  }
						],
						"includes": {
						  "tweets": [
							{
							  "created_at": "2020-04-29T17:01:38.000Z",
							  "text": "During these unprecedented times, what’s happening on Twitter can help the world better understand &amp; respond to the pandemic. \n\nWe're launching a free COVID-19 stream endpoint so qualified devs &amp; researchers can study the public conversation in real-time. https://t.co/BPqMcQzhId",
							  "id": "1255542774432063488"
							},
							{
							  "created_at": "2020-06-19T21:12:30.000Z",
							  "text": "Minneapolis\n @FredTJoseph https://t.co/lNTOkyguG1",
							  "id": "1274087687469715457"
							}
						  ]
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				ids: []string{"2244994945", "783214"},
				fieldOpts: UserFieldOptions{
					Expansions:  []Expansion{ExpansionPinnedTweetID},
					UserFields:  []UserField{UserFieldCreatedAt},
					TweetFields: []TweetField{TweetFieldCreatedAt},
				},
			},
			want: UserLookups{
				"2244994945": UserLookup{
					User: UserObj{
						UserName:      "TwitterDev",
						ID:            "2244994945",
						Name:          "Twitter Dev",
						CreatedAt:     "2013-12-14T04:35:55.000Z",
						PinnedTweetID: "1255542774432063488",
					},
					Tweet: &TweetObj{
						ID:        "1255542774432063488",
						Text:      "During these unprecedented times, what’s happening on Twitter can help the world better understand &amp; respond to the pandemic. \n\nWe're launching a free COVID-19 stream endpoint so qualified devs &amp; researchers can study the public conversation in real-time. https://t.co/BPqMcQzhId",
						CreatedAt: "2020-04-29T17:01:38.000Z",
					},
				},
				"783214": UserLookup{
					User: UserObj{
						UserName:      "Twitter",
						ID:            "783214",
						Name:          "Twitter",
						CreatedAt:     "2007-02-20T14:35:54.000Z",
						PinnedTweetID: "1274087687469715457",
					},
					Tweet: &TweetObj{
						ID:        "1274087687469715457",
						Text:      "Minneapolis\n @FredTJoseph https://t.co/lNTOkyguG1",
						CreatedAt: "2020-06-19T21:12:30.000Z",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "tweet error",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), userLookupEndpoint) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetLookupEndpoint)
					}
					body := `{
						"title": "Invalid Request",
						"detail": "One or more parameters to your request was invalid.",
						"type": "https://api.twitter.com/2/problems/invalid-request"
					}`
					return &http.Response{
						StatusCode: http.StatusBadRequest,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				ids: []string{"1067094924124872705"},
				fieldOpts: UserFieldOptions{
					UserFields: []UserField{UserFieldVerified, UserFieldUserName, UserFieldID, UserFieldName},
				},
			},
			want:    nil,
			wantErr: true,
			wantTweetError: &TweetErrorResponse{
				StatusCode: http.StatusBadRequest,
				Title:      "Invalid Request",
				Detail:     "One or more parameters to your request was invalid.",
				Type:       "https://api.twitter.com/2/problems/invalid-request",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := u.Lookup(context.Background(), tt.args.ids, tt.args.fieldOpts)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Lookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User.Lookup() = %v, want %v", got, tt.want)
			}
			var tweetErr *TweetErrorResponse
			if errors.As(err, &tweetErr) {
				if !reflect.DeepEqual(tweetErr, tt.wantTweetError) {
					t.Errorf("User.Lookup() Error = %+v, want %+v", tweetErr, tt.wantTweetError)
				}
			}
		})
	}
}

func TestUser_LookupUsername(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		usernames []string
		fieldOpts UserFieldOptions
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           UserLookups
		wantErr        bool
		wantTweetError *TweetErrorResponse
	}{
		{
			name: "success user name",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), userNameLookupEndpoint) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetLookupEndpoint)
					}
					body := `{
						"data": {
						  "id": "2244994945",
						  "name": "Twitter Dev",
						  "username": "TwitterDev"
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				usernames: []string{"TwitterDev"},
				fieldOpts: UserFieldOptions{},
			},
			want: UserLookups{
				"2244994945": UserLookup{
					User: UserObj{
						UserName: "TwitterDev",
						ID:       "2244994945",
						Name:     "Twitter Dev",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "success ids",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), userNamesLookupEndpoint) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetLookupEndpoint)
					}
					body := `{
						"data": [
						  {
							"created_at": "2013-12-14T04:35:55.000Z",
							"username": "TwitterDev",
							"pinned_tweet_id": "1255542774432063488",
							"id": "2244994945",
							"name": "Twitter Dev"
						  },
						  {
							"created_at": "2007-02-20T14:35:54.000Z",
							"username": "Twitter",
							"pinned_tweet_id": "1274087687469715457",
							"id": "783214",
							"name": "Twitter"
						  }
						],
						"includes": {
						  "tweets": [
							{
							  "created_at": "2020-04-29T17:01:38.000Z",
							  "text": "During these unprecedented times, what’s happening on Twitter can help the world better understand &amp; respond to the pandemic. \n\nWe're launching a free COVID-19 stream endpoint so qualified devs &amp; researchers can study the public conversation in real-time. https://t.co/BPqMcQzhId",
							  "id": "1255542774432063488"
							},
							{
							  "created_at": "2020-06-19T21:12:30.000Z",
							  "text": "Minneapolis\n @FredTJoseph https://t.co/lNTOkyguG1",
							  "id": "1274087687469715457"
							}
						  ]
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				usernames: []string{"TwitterDev", "Twitter"},
				fieldOpts: UserFieldOptions{
					Expansions:  []Expansion{ExpansionPinnedTweetID},
					UserFields:  []UserField{UserFieldCreatedAt},
					TweetFields: []TweetField{TweetFieldCreatedAt},
				},
			},
			want: UserLookups{
				"2244994945": UserLookup{
					User: UserObj{
						UserName:      "TwitterDev",
						ID:            "2244994945",
						Name:          "Twitter Dev",
						CreatedAt:     "2013-12-14T04:35:55.000Z",
						PinnedTweetID: "1255542774432063488",
					},
					Tweet: &TweetObj{
						ID:        "1255542774432063488",
						Text:      "During these unprecedented times, what’s happening on Twitter can help the world better understand &amp; respond to the pandemic. \n\nWe're launching a free COVID-19 stream endpoint so qualified devs &amp; researchers can study the public conversation in real-time. https://t.co/BPqMcQzhId",
						CreatedAt: "2020-04-29T17:01:38.000Z",
					},
				},
				"783214": UserLookup{
					User: UserObj{
						UserName:      "Twitter",
						ID:            "783214",
						Name:          "Twitter",
						CreatedAt:     "2007-02-20T14:35:54.000Z",
						PinnedTweetID: "1274087687469715457",
					},
					Tweet: &TweetObj{
						ID:        "1274087687469715457",
						Text:      "Minneapolis\n @FredTJoseph https://t.co/lNTOkyguG1",
						CreatedAt: "2020-06-19T21:12:30.000Z",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "tweet error",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), userNameLookupEndpoint) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetLookupEndpoint)
					}
					body := `{
						"title": "Invalid Request",
						"detail": "One or more parameters to your request was invalid.",
						"type": "https://api.twitter.com/2/problems/invalid-request"
					}`
					return &http.Response{
						StatusCode: http.StatusBadRequest,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				usernames: []string{"Twitter"},
				fieldOpts: UserFieldOptions{
					UserFields: []UserField{UserFieldVerified, UserFieldUserName, UserFieldID, UserFieldName},
				},
			},
			want:    nil,
			wantErr: true,
			wantTweetError: &TweetErrorResponse{
				StatusCode: http.StatusBadRequest,
				Title:      "Invalid Request",
				Detail:     "One or more parameters to your request was invalid.",
				Type:       "https://api.twitter.com/2/problems/invalid-request",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := u.LookupUsername(context.Background(), tt.args.usernames, tt.args.fieldOpts)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.LookupUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User.LookupUsername() = %v, want %v", got, tt.want)
			}
			var tweetErr *TweetErrorResponse
			if errors.As(err, &tweetErr) {
				if !reflect.DeepEqual(tweetErr, tt.wantTweetError) {
					t.Errorf("User.LookupUsername() Error = %+v, want %+v", tweetErr, tt.wantTweetError)
				}
			}
		})
	}
}

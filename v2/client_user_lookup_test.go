package twitter

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestClient_UserLookup(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		ids  []string
		opts UserLookupOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UserLookupResponse
		wantErr bool
	}{
		{
			name: "Success - Single ID Default",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), string(userLookupEndpoint)+"/2244994945") == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userLookupEndpoint)
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
				ids: []string{"2244994945"},
			},
			want: &UserLookupResponse{
				Raw: &UserRaw{
					Users: []*UserObj{
						{
							ID:       "2244994945",
							Name:     "Twitter Dev",
							UserName: "TwitterDev",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Success - Single ID Optional",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), string(userLookupEndpoint)+"/2244994945") == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userLookupEndpoint)
					}
					body := `{
						"data": {
						  "username": "TwitterDev",
						  "created_at": "2013-12-14T04:35:55.000Z",
						  "pinned_tweet_id": "1255542774432063488",
						  "id": "2244994945",
						  "name": "Twitter Dev"
						},
						"includes": {
						  "tweets": [
							{
							  "text": "During these unprecedented times, what’s happening on Twitter can help the world better understand &amp; respond to the pandemic. \n\nWe're launching a free COVID-19 stream endpoint so qualified devs &amp; researchers can study the public conversation in real-time. https://t.co/BPqMcQzhId",
							  "created_at": "2020-04-29T17:01:38.000Z",
							  "id": "1255542774432063488"
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
				ids: []string{"2244994945"},
				opts: UserLookupOpts{
					Expansions:  []Expansion{ExpansionPinnedTweetID},
					UserFields:  []UserField{UserFieldCreatedAt},
					TweetFields: []TweetField{TweetFieldCreatedAt},
				},
			},
			want: &UserLookupResponse{
				Raw: &UserRaw{
					Users: []*UserObj{
						{
							ID:            "2244994945",
							Name:          "Twitter Dev",
							UserName:      "TwitterDev",
							CreatedAt:     "2013-12-14T04:35:55.000Z",
							PinnedTweetID: "1255542774432063488",
						},
					},
					Includes: &UserRawIncludes{
						Tweets: []*TweetObj{
							{
								ID:        "1255542774432063488",
								CreatedAt: "2020-04-29T17:01:38.000Z",
								Text:      "During these unprecedented times, what’s happening on Twitter can help the world better understand &amp; respond to the pandemic. \n\nWe're launching a free COVID-19 stream endpoint so qualified devs &amp; researchers can study the public conversation in real-time. https://t.co/BPqMcQzhId",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Success - Multiple ID Default",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), string(userLookupEndpoint)) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userLookupEndpoint)
					}
					body := `{
						"data": [
						  {
							"id": "2244994945",
							"username": "TwitterDev",
							"name": "Twitter Dev"
						  },
						  {
							"id": "783214",
							"username": "Twitter",
							"name": "Twitter"
						  }
						]
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				ids: []string{"2244994945", "783214"},
			},
			want: &UserLookupResponse{
				Raw: &UserRaw{
					Users: []*UserObj{
						{
							ID:       "2244994945",
							Name:     "Twitter Dev",
							UserName: "TwitterDev",
						},
						{
							ID:       "783214",
							Name:     "Twitter",
							UserName: "Twitter",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Success - Multiple ID Optional",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), string(userLookupEndpoint)) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userLookupEndpoint)
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
							  "text": "📍 Minneapolis\n🗣️ @FredTJoseph https://t.co/lNTOkyguG1",
							  "id": "1274087687469715457"
							}
						  ]
						}
					  }
					  `
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				ids: []string{"2244994945", "783214"},
				opts: UserLookupOpts{
					Expansions:  []Expansion{ExpansionPinnedTweetID},
					UserFields:  []UserField{UserFieldCreatedAt},
					TweetFields: []TweetField{TweetFieldCreatedAt},
				},
			},
			want: &UserLookupResponse{
				Raw: &UserRaw{
					Users: []*UserObj{
						{
							ID:            "2244994945",
							Name:          "Twitter Dev",
							UserName:      "TwitterDev",
							CreatedAt:     "2013-12-14T04:35:55.000Z",
							PinnedTweetID: "1255542774432063488",
						},
						{
							ID:            "783214",
							Name:          "Twitter",
							UserName:      "Twitter",
							CreatedAt:     "2007-02-20T14:35:54.000Z",
							PinnedTweetID: "1274087687469715457",
						},
					},
					Includes: &UserRawIncludes{
						Tweets: []*TweetObj{
							{
								ID:        "1255542774432063488",
								CreatedAt: "2020-04-29T17:01:38.000Z",
								Text:      "During these unprecedented times, what’s happening on Twitter can help the world better understand &amp; respond to the pandemic. \n\nWe're launching a free COVID-19 stream endpoint so qualified devs &amp; researchers can study the public conversation in real-time. https://t.co/BPqMcQzhId",
							},
							{
								ID:        "1274087687469715457",
								CreatedAt: "2020-06-19T21:12:30.000Z",
								Text:      "📍 Minneapolis\n🗣️ @FredTJoseph https://t.co/lNTOkyguG1",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := c.UserLookup(context.Background(), tt.args.ids, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UserLookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.UserLookup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_UserNameLookup(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		usernames []string
		opts      UserLookupOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UserLookupResponse
		wantErr bool
	}{
		{
			name: "Success - Single ID Default",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), string(userNameLookupEndpoint)+"/username") == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userNameLookupEndpoint)
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
			},
			want: &UserLookupResponse{
				Raw: &UserRaw{
					Users: []*UserObj{
						{
							ID:       "2244994945",
							Name:     "Twitter Dev",
							UserName: "TwitterDev",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Success - Single ID Optional",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), string(userNameLookupEndpoint)+"/username") == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userNameLookupEndpoint)
					}
					body := `{
						"data": {
						  "username": "TwitterDev",
						  "created_at": "2013-12-14T04:35:55.000Z",
						  "pinned_tweet_id": "1255542774432063488",
						  "id": "2244994945",
						  "name": "Twitter Dev"
						},
						"includes": {
						  "tweets": [
							{
							  "text": "During these unprecedented times, what’s happening on Twitter can help the world better understand &amp; respond to the pandemic. \n\nWe're launching a free COVID-19 stream endpoint so qualified devs &amp; researchers can study the public conversation in real-time. https://t.co/BPqMcQzhId",
							  "created_at": "2020-04-29T17:01:38.000Z",
							  "id": "1255542774432063488"
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
				usernames: []string{"TwitterDev"},
				opts: UserLookupOpts{
					Expansions:  []Expansion{ExpansionPinnedTweetID},
					UserFields:  []UserField{UserFieldCreatedAt},
					TweetFields: []TweetField{TweetFieldCreatedAt},
				},
			},
			want: &UserLookupResponse{
				Raw: &UserRaw{
					Users: []*UserObj{
						{
							ID:            "2244994945",
							Name:          "Twitter Dev",
							UserName:      "TwitterDev",
							CreatedAt:     "2013-12-14T04:35:55.000Z",
							PinnedTweetID: "1255542774432063488",
						},
					},
					Includes: &UserRawIncludes{
						Tweets: []*TweetObj{
							{
								ID:        "1255542774432063488",
								CreatedAt: "2020-04-29T17:01:38.000Z",
								Text:      "During these unprecedented times, what’s happening on Twitter can help the world better understand &amp; respond to the pandemic. \n\nWe're launching a free COVID-19 stream endpoint so qualified devs &amp; researchers can study the public conversation in real-time. https://t.co/BPqMcQzhId",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Success - Multiple ID Default",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), string(userNameLookupEndpoint)) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userNameLookupEndpoint)
					}
					body := `{
						"data": [
						  {
							"id": "2244994945",
							"username": "TwitterDev",
							"name": "Twitter Dev"
						  },
						  {
							"id": "783214",
							"username": "Twitter",
							"name": "Twitter"
						  }
						]
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				usernames: []string{"TwitterDev", "Twitter"},
			},
			want: &UserLookupResponse{
				Raw: &UserRaw{
					Users: []*UserObj{
						{
							ID:       "2244994945",
							Name:     "Twitter Dev",
							UserName: "TwitterDev",
						},
						{
							ID:       "783214",
							Name:     "Twitter",
							UserName: "Twitter",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Success - Multiple ID Optional",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), string(userNameLookupEndpoint)) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userNameLookupEndpoint)
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
							  "text": "📍 Minneapolis\n🗣️ @FredTJoseph https://t.co/lNTOkyguG1",
							  "id": "1274087687469715457"
							}
						  ]
						}
					  }
					  `
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				usernames: []string{"TwitterDev", "Twitter"},
				opts: UserLookupOpts{
					Expansions:  []Expansion{ExpansionPinnedTweetID},
					UserFields:  []UserField{UserFieldCreatedAt},
					TweetFields: []TweetField{TweetFieldCreatedAt},
				},
			},
			want: &UserLookupResponse{
				Raw: &UserRaw{
					Users: []*UserObj{
						{
							ID:            "2244994945",
							Name:          "Twitter Dev",
							UserName:      "TwitterDev",
							CreatedAt:     "2013-12-14T04:35:55.000Z",
							PinnedTweetID: "1255542774432063488",
						},
						{
							ID:            "783214",
							Name:          "Twitter",
							UserName:      "Twitter",
							CreatedAt:     "2007-02-20T14:35:54.000Z",
							PinnedTweetID: "1274087687469715457",
						},
					},
					Includes: &UserRawIncludes{
						Tweets: []*TweetObj{
							{
								ID:        "1255542774432063488",
								CreatedAt: "2020-04-29T17:01:38.000Z",
								Text:      "During these unprecedented times, what’s happening on Twitter can help the world better understand &amp; respond to the pandemic. \n\nWe're launching a free COVID-19 stream endpoint so qualified devs &amp; researchers can study the public conversation in real-time. https://t.co/BPqMcQzhId",
							},
							{
								ID:        "1274087687469715457",
								CreatedAt: "2020-06-19T21:12:30.000Z",
								Text:      "📍 Minneapolis\n🗣️ @FredTJoseph https://t.co/lNTOkyguG1",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := c.UserNameLookup(context.Background(), tt.args.usernames, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UserNameLookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.UserNameLookup() = %v, want %v", got, tt.want)
			}
		})
	}
}

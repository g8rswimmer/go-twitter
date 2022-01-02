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

func TestClient_UserRetweet(t *testing.T) {
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
		want    *UserRetweetResponse
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodPost {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodPost)
					}
					if strings.Contains(req.URL.String(), userManageRetweetEndpoint.urlID("", "user-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userManageRetweetEndpoint)
					}
					body := `{
						"data": {
						  "retweeted": true
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				userID:  "user-1234",
				tweetID: "tweet-1234",
			},
			want: &UserRetweetResponse{
				Data: &RetweetData{
					Retweeted: true,
				},
			},
			wantErr: false,
		},
		{
			name:   "No User ID",
			fields: fields{},
			args: args{
				tweetID: "tweet-1234",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "No Tweet ID",
			fields: fields{},
			args: args{
				userID: "user-1234",
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
			got, err := c.UserRetweet(context.Background(), tt.args.userID, tt.args.tweetID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UserRetweet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.UserRetweet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_DeleteUserRetweet(t *testing.T) {
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
		want    *DeleteUserRetweetResponse
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodDelete {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodDelete)
					}
					if strings.Contains(req.URL.String(), userManageRetweetEndpoint.urlID("", "user-1234")+"/tweet-1234") == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userManageRetweetEndpoint)
					}
					body := `{
							"data": {
							  "retweeted": false
							}
						  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				userID:  "user-1234",
				tweetID: "tweet-1234",
			},
			want: &DeleteUserRetweetResponse{
				Data: &RetweetData{
					Retweeted: false,
				},
			},
			wantErr: false,
		},
		{
			name:   "No User ID",
			fields: fields{},
			args: args{
				tweetID: "tweet-1234",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "No Tweet ID",
			fields: fields{},
			args: args{
				userID: "user-1234",
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
			got, err := c.DeleteUserRetweet(context.Background(), tt.args.userID, tt.args.tweetID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteUserRetweet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.DeleteUserRetweet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_UserRetweetLookup(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		tweetID string
		opts    UserRetweetLookuoOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UserRetweetLookupResponse
		wantErr bool
	}{
		{
			name: "Success - no options",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), userRetweetLookupEndpoint.urlID("", "tweet-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userRetweetLookupEndpoint)
					}
					body := `{
						"data": [
						  {
							"id": "1065249714214457345",
							"name": "Spaces",
							"username": "TwitterSpaces"
						  },
						  {
							"id": "783214",
							"name": "Twitter",
							"username": "Twitter"
						  },
						  {
							"id": "1526228120",
							"name": "Twitter Data",
							"username": "TwitterData"
						  }
						],
						"meta": {
						  "result_count": 3
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				tweetID: "tweet-1234",
			},
			want: &UserRetweetLookupResponse{
				Raw: &UserRetweetRaw{
					Users: []*UserObj{
						{
							ID:       "1065249714214457345",
							Name:     "Spaces",
							UserName: "TwitterSpaces",
						},
						{
							ID:       "783214",
							Name:     "Twitter",
							UserName: "Twitter",
						},
						{
							ID:       "1526228120",
							Name:     "Twitter Data",
							UserName: "TwitterData",
						},
					},
				},
				Meta: &UserRetweetMeta{
					ResultCount: 3,
				},
			},
			wantErr: false,
		},
		{
			name: "Success - with options",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), userRetweetLookupEndpoint.urlID("", "tweet-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userRetweetLookupEndpoint)
					}
					body := `{
						"data": [
						  {
							"id": "1065249714214457345",
							"created_at": "2018-11-21T14:24:58.000Z",
							"name": "Spaces",
							"pinned_tweet_id": "1389270063807598594",
							"description": "Twitter Spaces is where live audio conversations happen.",
							"username": "TwitterSpaces"
						  },
						  {
							"id": "783214",
							"created_at": "2007-02-20T14:35:54.000Z",
							"name": "Twitter",
							"description": "What's happening?!",
							"username": "Twitter"
						  },
						  {
							"id": "1526228120",
							"created_at": "2013-06-17T23:57:45.000Z",
							"name": "Twitter Data",
							"description": "Data-driven insights about notable moments and conversations from Twitter, Inc., plus tips and tricks to help you get the most out of Twitter data.",
							"username": "TwitterData"
						  },
						  {
							"id": "2244994945",
							"created_at": "2013-12-14T04:35:55.000Z",
							"name": "Twitter Dev",
							"pinned_tweet_id": "1354143047324299264",
							"description": "The voice of the #TwitterDev team and your official source for updates, news, and events, related to the #TwitterAPI.",
							"username": "TwitterDev"
						  },
						  {
							"id": "6253282",
							"created_at": "2007-05-23T06:01:13.000Z",
							"name": "Twitter API",
							"pinned_tweet_id": "1293595870563381249",
							"description": "Tweets about changes and service issues. Follow @TwitterDev for more.",
							"username": "TwitterAPI"
						  }
						],
						"includes": {
							"tweets": [
							  {
								"id": "1389270063807598594",
								"text": "now, everyone with 600 or more followers can host a Space.nnbased on what we've learned, these accounts are likely to have a good experience hosting because of their existing audience. before bringing the ability to create a Space to everyone, we're focused on a few things. :thread:"
							  },
							  {
								"id": "1354143047324299264",
								"text": "Academics are one of the biggest groups using the #TwitterAPI to research what's happening. Their work helps make the world (&amp; Twitter) a better place, and now more than ever, we must enable more of it. nIntroducing :drum_with_drumsticks: the Academic Research product track!nhttps://t.co/nOFiGewAV2"
							  },
							  {
								"id": "1293595870563381249",
								"text": "Twitter API v2: Early Access releasednnToday we announced Early Access to the first endpoints of the new Twitter API!nn#TwitterAPI #EarlyAccess #VersionBump https://t.co/g7v3aeIbtQ"
							  }
							]
						},
						"meta": {
							"result_count": 5
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				tweetID: "tweet-1234",
				opts: UserRetweetLookuoOpts{
					Expansions:  []Expansion{ExpansionPinnedTweetID},
					UserFields:  []UserField{UserFieldCreatedAt, UserFieldDescription},
					TweetFields: []TweetField{TweetFieldCreatedAt},
				},
			},
			want: &UserRetweetLookupResponse{
				Raw: &UserRetweetRaw{
					Users: []*UserObj{
						{
							ID:            "1065249714214457345",
							Name:          "Spaces",
							UserName:      "TwitterSpaces",
							CreatedAt:     "2018-11-21T14:24:58.000Z",
							PinnedTweetID: "1389270063807598594",
							Description:   "Twitter Spaces is where live audio conversations happen.",
						},
						{
							ID:          "783214",
							Name:        "Twitter",
							UserName:    "Twitter",
							CreatedAt:   "2007-02-20T14:35:54.000Z",
							Description: "What's happening?!",
						},
						{
							ID:          "1526228120",
							Name:        "Twitter Data",
							UserName:    "TwitterData",
							CreatedAt:   "2013-06-17T23:57:45.000Z",
							Description: "Data-driven insights about notable moments and conversations from Twitter, Inc., plus tips and tricks to help you get the most out of Twitter data.",
						},
						{
							ID:            "2244994945",
							Name:          "Twitter Dev",
							UserName:      "TwitterDev",
							CreatedAt:     "2013-12-14T04:35:55.000Z",
							PinnedTweetID: "1354143047324299264",
							Description:   "The voice of the #TwitterDev team and your official source for updates, news, and events, related to the #TwitterAPI.",
						},
						{
							ID:            "6253282",
							Name:          "Twitter API",
							UserName:      "TwitterAPI",
							CreatedAt:     "2007-05-23T06:01:13.000Z",
							PinnedTweetID: "1293595870563381249",
							Description:   "Tweets about changes and service issues. Follow @TwitterDev for more.",
						},
					},
					Includes: &UserRetweetRawIncludes{
						Tweets: []*TweetObj{
							{
								ID:   "1389270063807598594",
								Text: "now, everyone with 600 or more followers can host a Space.nnbased on what we've learned, these accounts are likely to have a good experience hosting because of their existing audience. before bringing the ability to create a Space to everyone, we're focused on a few things. :thread:",
							},
							{
								ID:   "1354143047324299264",
								Text: "Academics are one of the biggest groups using the #TwitterAPI to research what's happening. Their work helps make the world (&amp; Twitter) a better place, and now more than ever, we must enable more of it. nIntroducing :drum_with_drumsticks: the Academic Research product track!nhttps://t.co/nOFiGewAV2",
							},
							{
								ID:   "1293595870563381249",
								Text: "Twitter API v2: Early Access releasednnToday we announced Early Access to the first endpoints of the new Twitter API!nn#TwitterAPI #EarlyAccess #VersionBump https://t.co/g7v3aeIbtQ",
							},
						},
					},
				},
				Meta: &UserRetweetMeta{
					ResultCount: 5,
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
			got, err := c.UserRetweetLookup(context.Background(), tt.args.tweetID, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UserRetweetLookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.UserRetweetLookup() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func TestClient_TweetLikesLookup(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		tweetID string
		opts    TweetLikesLookupOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TweetLikesLookupResponse
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
					if strings.Contains(req.URL.String(), tweetLikesEndpoint.urlID("", "tweet-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetLikesEndpoint)
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
							"result_count": 2,
							"next_token": "7140dibdnow9c7btw3w29grvxfcgvpb9n9coehpk7xz5i"
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
			want: &TweetLikesLookupResponse{
				Raw: &UserRaw{
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
				Meta: &TweetLikesMeta{
					ResultCount: 2,
					NextToken:   "7140dibdnow9c7btw3w29grvxfcgvpb9n9coehpk7xz5i",
				},
				RateLimit: &RateLimit{
					Limit:     15,
					Remaining: 12,
					Reset:     Epoch(1644461060),
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
					if strings.Contains(req.URL.String(), tweetLikesEndpoint.urlID("", "tweet-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetLikesEndpoint)
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
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				tweetID: "tweet-1234",
				opts: TweetLikesLookupOpts{
					Expansions:  []Expansion{ExpansionPinnedTweetID},
					UserFields:  []UserField{UserFieldCreatedAt, UserFieldDescription},
					TweetFields: []TweetField{TweetFieldCreatedAt},
				},
			},
			want: &TweetLikesLookupResponse{
				Raw: &UserRaw{
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
					Includes: &UserRawIncludes{
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
			got, err := c.TweetLikesLookup(context.Background(), tt.args.tweetID, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.TweetLikesLookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.TweetLikesLookup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_UserLikesLookup(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		userID string
		opts   UserLikesLookupOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UserLikesLookupResponse
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
					if strings.Contains(req.URL.String(), userLikedTweetEndpoint.urlID("", "2244994945")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userLikedTweetEndpoint)
					}
					body := `{
						"data": [
						  {
							"id": "1338971066773905408",
							"text": "💡 Using Twitter data for academic research? Join our next livestream this Friday @ 9am PT on https://t.co/GrtBOXh5Y1!\n \n@SuhemParack will show how to get started with recent search &amp; filtered stream endpoints on the #TwitterAPI v2, the new Tweet payload, annotations, &amp; more. https://t.co/IraD2Z7wEg"
						  },
						  {
							"id": "1338923691497959425",
							"text": "📈 Live now with @jessicagarson and @i_am_daniele! https://t.co/Y1AFzsTTxb"
						  }
						],
						"meta": {
						  "result_count": 2,
						  "next_token": "7140dibdnow9c7btw3w29grvxfcgvpb9n9coehpk7xz5i"
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
				userID: "2244994945",
			},
			want: &UserLikesLookupResponse{
				Raw: &TweetRaw{
					Tweets: []*TweetObj{
						{
							ID:   "1338971066773905408",
							Text: "💡 Using Twitter data for academic research? Join our next livestream this Friday @ 9am PT on https://t.co/GrtBOXh5Y1!\n \n@SuhemParack will show how to get started with recent search &amp; filtered stream endpoints on the #TwitterAPI v2, the new Tweet payload, annotations, &amp; more. https://t.co/IraD2Z7wEg",
						},
						{
							ID:   "1338923691497959425",
							Text: "📈 Live now with @jessicagarson and @i_am_daniele! https://t.co/Y1AFzsTTxb",
						},
					},
				},
				Meta: &UserLikesMeta{
					ResultCount: 2,
					NextToken:   "7140dibdnow9c7btw3w29grvxfcgvpb9n9coehpk7xz5i",
				},
				RateLimit: &RateLimit{
					Limit:     15,
					Remaining: 12,
					Reset:     Epoch(1644461060),
				},
			},
			wantErr: false,
		},
		{
			name: "Success - Optional Fields",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), userLikedTweetEndpoint.urlID("", "2244994945")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userLikedTweetEndpoint)
					}
					body := `{
						"data": [
						  {
							"author_id": "2244994945",
							"conversation_id": "1338971066773905408",
							"text": "💡 Using Twitter data for academic research? Join our next livestream this Friday @ 9am PT on https://t.co/GrtBOXh5Y1!\n \n@SuhemParack will show how to get started with recent search &amp; filtered stream endpoints on the #TwitterAPI v2, the new Tweet payload, annotations, &amp; more. https://t.co/IraD2Z7wEg",
							"context_annotations": [
							  {
								"domain": {
								  "id": "47",
								  "name": "Brand",
								  "description": "Brands and Companies"
								},
								"entity": {
								  "id": "10045225402",
								  "name": "Twitter"
								}
							  }
							],
							"public_metrics": {
							  "retweet_count": 10,
							  "reply_count": 1,
							  "like_count": 41,
							  "quote_count": 4
							},
							"id": "1338971066773905408",
							"created_at": "2020-12-15T22:15:53.000Z"
						  },
						  {
							"author_id": "2244994945",
							"conversation_id": "1338923691497959425",
							"text": "📈 Live now with @jessicagarson and @i_am_daniele! https://t.co/Y1AFzsTTxb",
							"context_annotations": [
							  {
								"domain": {
								  "id": "47",
								  "name": "Brand",
								  "description": "Brands and Companies"
								},
								"entity": {
								  "id": "10026378521",
								  "name": "Google "
								}
							  }
							],
							"public_metrics": {
							  "retweet_count": 3,
							  "reply_count": 0,
							  "like_count": 12,
							  "quote_count": 1
							},
							"id": "1338923691497959425",
							"created_at": "2020-12-15T19:07:38.000Z"
						  }
						],
						"includes": {
						  "users": [
							{
							  "id": "2244994945",
							  "name": "Twitter Dev",
							  "username": "TwitterDev"
							}
						  ]
						},
						"meta": {
						  "result_count": 2,
						  "next_token": "7140dibdnow9c7btw3w29n4v1mtag9kegr0gr7y26pnw3"
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				userID: "2244994945",
				opts: UserLikesLookupOpts{
					TweetFields: []TweetField{TweetFieldCreatedAt, TweetFieldAuthorID, TweetFieldConversationID, TweetFieldPublicMetrics, TweetFieldContextAnnotations},
					UserFields:  []UserField{UserFieldUserName},
					Expansions:  []Expansion{ExpansionAuthorID},
					MaxResults:  10,
				},
			},
			want: &UserLikesLookupResponse{
				Raw: &TweetRaw{
					Tweets: []*TweetObj{
						{
							ID:             "1338971066773905408",
							Text:           "💡 Using Twitter data for academic research? Join our next livestream this Friday @ 9am PT on https://t.co/GrtBOXh5Y1!\n \n@SuhemParack will show how to get started with recent search &amp; filtered stream endpoints on the #TwitterAPI v2, the new Tweet payload, annotations, &amp; more. https://t.co/IraD2Z7wEg",
							AuthorID:       "2244994945",
							ConversationID: "1338971066773905408",
							CreatedAt:      "2020-12-15T22:15:53.000Z",
							ContextAnnotations: []*TweetContextAnnotationObj{
								{
									Domain: TweetContextObj{
										ID:          "47",
										Name:        "Brand",
										Description: "Brands and Companies",
									},
									Entity: TweetContextObj{
										ID:   "10045225402",
										Name: "Twitter",
									},
								},
							},
							PublicMetrics: &TweetMetricsObj{
								Retweets: 10,
								Replies:  1,
								Likes:    41,
								Quotes:   4,
							},
						},
						{
							ID:             "1338923691497959425",
							Text:           "📈 Live now with @jessicagarson and @i_am_daniele! https://t.co/Y1AFzsTTxb",
							AuthorID:       "2244994945",
							ConversationID: "1338923691497959425",
							CreatedAt:      "2020-12-15T19:07:38.000Z",
							ContextAnnotations: []*TweetContextAnnotationObj{
								{
									Domain: TweetContextObj{
										ID:          "47",
										Name:        "Brand",
										Description: "Brands and Companies",
									},
									Entity: TweetContextObj{
										ID:   "10026378521",
										Name: "Google ",
									},
								},
							},
							PublicMetrics: &TweetMetricsObj{
								Retweets: 3,
								Replies:  0,
								Likes:    12,
								Quotes:   1,
							},
						},
					},
					Includes: &TweetRawIncludes{
						Users: []*UserObj{
							{
								ID:       "2244994945",
								Name:     "Twitter Dev",
								UserName: "TwitterDev",
							},
						},
					},
				},
				Meta: &UserLikesMeta{
					ResultCount: 2,
					NextToken:   "7140dibdnow9c7btw3w29n4v1mtag9kegr0gr7y26pnw3",
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
			got, err := c.UserLikesLookup(context.Background(), tt.args.userID, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UserLikesLookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.UserLikesLookup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_UserLikes(t *testing.T) {
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
		want    *UserLikesResponse
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
					if strings.Contains(req.URL.String(), userLikesEndpoint.urlID("", "user-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userLikesEndpoint)
					}
					body := `{
						"data": {
						  "liked": true
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
				tweetID: "tweet-1234",
			},
			want: &UserLikesResponse{
				Data: &UserLikesData{
					Liked: true,
				},
				RateLimit: &RateLimit{
					Limit:     15,
					Remaining: 12,
					Reset:     Epoch(1644461060),
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
			got, err := c.UserLikes(context.Background(), tt.args.userID, tt.args.tweetID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UserLikes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.UserLikes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_DeleteUserLikes(t *testing.T) {
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
		want    *DeleteUserLikesResponse
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
					if strings.Contains(req.URL.String(), userLikesEndpoint.urlID("", "user-1234")+"/tweet-1234") == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userLikesEndpoint)
					}
					body := `{
							"data": {
							  "liked": false
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
				tweetID: "tweet-1234",
			},
			want: &DeleteUserLikesResponse{
				Data: &UserLikesData{
					Liked: false,
				},
				RateLimit: &RateLimit{
					Limit:     15,
					Remaining: 12,
					Reset:     Epoch(1644461060),
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
			got, err := c.DeleteUserLikes(context.Background(), tt.args.userID, tt.args.tweetID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteUserLikes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.DeleteUserLikes() = %v, want %v", got, tt.want)
			}
		})
	}
}

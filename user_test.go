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

func TestUser_LookupFollowing(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		id         string
		followOpts UserFollowOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UserFollowLookup
		wantErr bool
	}{
		{
			name: "Success-Basic",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), "2/users/2244994945/following?max_results=10") == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetLookupEndpoint)
					}
					body := `{
						"data": [
						  {
							"id": "6253282",
							"name": "Twitter API",
							"username": "TwitterAPI"
						  },
						  {
							"id": "2244994945",
							"name": "Twitter Dev",
							"username": "TwitterDev"
						  },
						  {
							"id": "783214",
							"name": "Twitter",
							"username": "Twitter"
						  },
						  {
							"id": "95731075",
							"name": "Twitter Safety",
							"username": "TwitterSafety"
						  },
						  {
							"id": "3260518932",
							"name": "Twitter Moments",
							"username": "TwitterMoments"
						  },
						  {
							"id": "373471064",
							"name": "Twitter Music",
							"username": "TwitterMusic"
						  },
						  {
							"id": "791978718",
							"name": "Twitter Official Partner",
							"username": "OfficialPartner"
						  },
						  {
							"id": "17874544",
							"name": "Twitter Support",
							"username": "TwitterSupport"
						  },
						  {
							"id": "234489024",
							"name": "Twitter Comms",
							"username": "TwitterComms"
						  },
						  {
							"id": "1526228120",
							"name": "Twitter Data",
							"username": "TwitterData"
						  }
						],
						"meta": {
						  "result_count": 10,
						  "next_token": "DFEDBNRFT3MHCZZZ"
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				id: "2244994945",
				followOpts: UserFollowOptions{
					MaxResults: 10,
				},
			},
			want: &UserFollowLookup{
				Lookups: UserLookups{
					"6253282": UserLookup{
						User: UserObj{
							ID:       "6253282",
							Name:     "Twitter API",
							UserName: "TwitterAPI",
						},
					},
					"2244994945": UserLookup{
						User: UserObj{
							ID:       "2244994945",
							Name:     "Twitter Dev",
							UserName: "TwitterDev",
						},
					},
					"783214": UserLookup{
						User: UserObj{
							ID:       "783214",
							Name:     "Twitter",
							UserName: "Twitter",
						},
					},
					"95731075": UserLookup{
						User: UserObj{
							ID:       "95731075",
							Name:     "Twitter Safety",
							UserName: "TwitterSafety",
						},
					},
					"3260518932": UserLookup{
						User: UserObj{
							ID:       "3260518932",
							Name:     "Twitter Moments",
							UserName: "TwitterMoments",
						},
					},
					"373471064": UserLookup{
						User: UserObj{
							ID:       "373471064",
							Name:     "Twitter Music",
							UserName: "TwitterMusic",
						},
					},
					"791978718": UserLookup{
						User: UserObj{
							ID:       "791978718",
							Name:     "Twitter Official Partner",
							UserName: "OfficialPartner",
						},
					},
					"17874544": UserLookup{
						User: UserObj{
							ID:       "17874544",
							Name:     "Twitter Support",
							UserName: "TwitterSupport",
						},
					},
					"234489024": UserLookup{
						User: UserObj{
							ID:       "234489024",
							Name:     "Twitter Comms",
							UserName: "TwitterComms",
						},
					},
					"1526228120": UserLookup{
						User: UserObj{
							ID:       "1526228120",
							Name:     "Twitter Data",
							UserName: "TwitterData",
						},
					},
				},
				Meta: &UserFollowMeta{
					ResultCount: 10,
					NextToken:   "DFEDBNRFT3MHCZZZ",
				},
			},
			wantErr: false,
		},
		{
			name: "Success-Optional",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), "2/users/2244994945/following?") == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetLookupEndpoint)
					}
					body := `{
						"data": [
						  {
							"pinned_tweet_id": "1293595870563381249",
							"id": "6253282",
							"username": "TwitterAPI",
							"name": "Twitter API"
						  },
						  {
							"pinned_tweet_id": "1293593516040269825",
							"id": "2244994945",
							"username": "TwitterDev",
							"name": "Twitter Dev"
						  },
						  {
							"id": "783214",
							"username": "Twitter",
							"name": "Twitter"
						  },
						  {
							"pinned_tweet_id": "1271186240323432452",
							"id": "95731075",
							"username": "TwitterSafety",
							"name": "Twitter Safety"
						  },
						  {
							"id": "3260518932",
							"username": "TwitterMoments",
							"name": "Twitter Moments"
						  },
						  {
							"pinned_tweet_id": "1293216056274759680",
							"id": "373471064",
							"username": "TwitterMusic",
							"name": "Twitter Music"
						  },
						  {
							"id": "791978718",
							"username": "OfficialPartner",
							"name": "Twitter Official Partner"
						  },
						  {
							"pinned_tweet_id": "1289000334497439744",
							"id": "17874544",
							"username": "TwitterSupport",
							"name": "Twitter Support"
						  },
						  {
							"pinned_tweet_id": "1283543147444711424",
							"id": "234489024",
							"username": "TwitterComms",
							"name": "Twitter Comms"
						  },
						  {
							"id": "1526228120",
							"username": "TwitterData",
							"name": "Twitter Data"
						  }
						],
						"includes": {
						  "tweets": [
							{
							  "context_annotations": [
								{
								  "domain": {
									"id": "46",
									"name": "Brand Category",
									"description": "Categories within Brand Verticals that narrow down the scope of Brands"
								  },
								  "entity": {
									"id": "781974596752842752",
									"name": "Services"
								  }
								},
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
								},
								{
								  "domain": {
									"id": "65",
									"name": "Interests and Hobbies Vertical",
									"description": "Top level interests and hobbies groupings, like Food or Travel"
								  },
								  "entity": {
									"id": "848920371311001600",
									"name": "Technology",
									"description": "Technology and computing"
								  }
								},
								{
								  "domain": {
									"id": "66",
									"name": "Interests and Hobbies Category",
									"description": "A grouping of interests and hobbies entities, like Novelty Food or Destinations"
								  },
								  "entity": {
									"id": "848921413196984320",
									"name": "Computer programming",
									"description": "Computer programming"
								  }
								},
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
							  "id": "1293595870563381249",
							  "text": "Twitter API v2: Early Access released\n\nToday we announced Early Access to the first endpoints of the new Twitter API!\n\n#TwitterAPI #EarlyAccess #VersionBump https://t.co/g7v3aeIbtQ"
							},
							{
							  "context_annotations": [
								{
								  "domain": {
									"id": "46",
									"name": "Brand Category",
									"description": "Categories within Brand Verticals that narrow down the scope of Brands"
								  },
								  "entity": {
									"id": "781974596752842752",
									"name": "Services"
								  }
								},
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
								},
								{
								  "domain": {
									"id": "65",
									"name": "Interests and Hobbies Vertical",
									"description": "Top level interests and hobbies groupings, like Food or Travel"
								  },
								  "entity": {
									"id": "848920371311001600",
									"name": "Technology",
									"description": "Technology and computing"
								  }
								},
								{
								  "domain": {
									"id": "66",
									"name": "Interests and Hobbies Category",
									"description": "A grouping of interests and hobbies entities, like Novelty Food or Destinations"
								  },
								  "entity": {
									"id": "848921413196984320",
									"name": "Computer programming",
									"description": "Computer programming"
								  }
								}
							  ],
							  "id": "1293593516040269825",
							  "text": "It’s finally here! 🥁 Say hello to the new #TwitterAPI.\n\nWe’re rebuilding the Twitter API v2 from the ground up to better serve our developer community. And today’s launch is only the beginning.\n\nhttps://t.co/32VrwpGaJw https://t.co/KaFSbjWUA8"
							},
							{
							  "id": "1271186240323432452",
							  "text": "We’re disclosing new state-linked information operations to our public archive — the only one of its kind in the industry. Originating from the People’s Republic of China (PRC), Russia, and Turkey, all associated accounts and content have been removed. https://t.co/obRqr96iYm"
							},
							{
							  "id": "1293216056274759680",
							  "text": "say howdy to your new yeehaw king @orvillepeck—our #ArtistToFollow this month 🤠 https://t.co/3pk9fYcPHb"
							},
							{
							  "context_annotations": [
								{
								  "domain": {
									"id": "46",
									"name": "Brand Category",
									"description": "Categories within Brand Verticals that narrow down the scope of Brands"
								  },
								  "entity": {
									"id": "781974596752842752",
									"name": "Services"
								  }
								},
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
							  "id": "1289000334497439744",
							  "text": "We’ve significantly limited access to our internal tools and systems. Until we can safely resume normal operations, our response times to some support needs and reports will be slower. Thank you for your patience as we work through this."
							},
							{
							  "context_annotations": [
								{
								  "domain": {
									"id": "46",
									"name": "Brand Category",
									"description": "Categories within Brand Verticals that narrow down the scope of Brands"
								  },
								  "entity": {
									"id": "781974596752842752",
									"name": "Services"
								  }
								},
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
							  "id": "1283543147444711424",
							  "text": "Follow @TwitterSupport for the latest on the security incident ⬇️ https://t.co/7FKKksJqxV"
							}
						  ]
						},
						"meta": {
							"result_count": 10,
							"next_token": "DFEDBNRFT3MHCZZZ"
					    }
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				id: "2244994945",
				followOpts: UserFollowOptions{
					MaxResults:  10,
					Expansions:  []Expansion{ExpansionPinnedTweetID},
					TweetFields: []TweetField{TweetFieldID, TweetFieldContextAnnotations},
				},
			},
			want: &UserFollowLookup{
				Lookups: UserLookups{
					"6253282": UserLookup{
						User: UserObj{
							ID:            "6253282",
							Name:          "Twitter API",
							UserName:      "TwitterAPI",
							PinnedTweetID: "1293595870563381249",
						},
						Tweet: &TweetObj{
							ID:   "1293595870563381249",
							Text: "Twitter API v2: Early Access released\n\nToday we announced Early Access to the first endpoints of the new Twitter API!\n\n#TwitterAPI #EarlyAccess #VersionBump https://t.co/g7v3aeIbtQ",
							ContextAnnotations: []TweetContextAnnotationObj{
								{
									Domain: TweetContextObj{
										ID:          "46",
										Name:        "Brand Category",
										Description: "Categories within Brand Verticals that narrow down the scope of Brands",
									},
									Entity: TweetContextObj{
										ID:   "781974596752842752",
										Name: "Services",
									},
								},
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
								{
									Domain: TweetContextObj{
										ID:          "65",
										Name:        "Interests and Hobbies Vertical",
										Description: "Top level interests and hobbies groupings, like Food or Travel",
									},
									Entity: TweetContextObj{
										ID:          "848920371311001600",
										Name:        "Technology",
										Description: "Technology and computing",
									},
								},
								{
									Domain: TweetContextObj{
										ID:          "66",
										Name:        "Interests and Hobbies Category",
										Description: "A grouping of interests and hobbies entities, like Novelty Food or Destinations",
									},
									Entity: TweetContextObj{
										ID:          "848921413196984320",
										Name:        "Computer programming",
										Description: "Computer programming",
									},
								},
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
						},
					},
					"2244994945": UserLookup{
						User: UserObj{
							ID:            "2244994945",
							Name:          "Twitter Dev",
							UserName:      "TwitterDev",
							PinnedTweetID: "1293593516040269825",
						},
						Tweet: &TweetObj{
							ID:   "1293593516040269825",
							Text: "It’s finally here! 🥁 Say hello to the new #TwitterAPI.\n\nWe’re rebuilding the Twitter API v2 from the ground up to better serve our developer community. And today’s launch is only the beginning.\n\nhttps://t.co/32VrwpGaJw https://t.co/KaFSbjWUA8",
							ContextAnnotations: []TweetContextAnnotationObj{
								{
									Domain: TweetContextObj{
										ID:          "46",
										Name:        "Brand Category",
										Description: "Categories within Brand Verticals that narrow down the scope of Brands",
									},
									Entity: TweetContextObj{
										ID:   "781974596752842752",
										Name: "Services",
									},
								},
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
								{
									Domain: TweetContextObj{
										ID:          "65",
										Name:        "Interests and Hobbies Vertical",
										Description: "Top level interests and hobbies groupings, like Food or Travel",
									},
									Entity: TweetContextObj{
										ID:          "848920371311001600",
										Name:        "Technology",
										Description: "Technology and computing",
									},
								},
								{
									Domain: TweetContextObj{
										ID:          "66",
										Name:        "Interests and Hobbies Category",
										Description: "A grouping of interests and hobbies entities, like Novelty Food or Destinations",
									},
									Entity: TweetContextObj{
										ID:          "848921413196984320",
										Name:        "Computer programming",
										Description: "Computer programming",
									},
								},
							},
						},
					},
					"783214": UserLookup{
						User: UserObj{
							ID:       "783214",
							Name:     "Twitter",
							UserName: "Twitter",
						},
					},
					"95731075": UserLookup{
						User: UserObj{
							ID:            "95731075",
							Name:          "Twitter Safety",
							UserName:      "TwitterSafety",
							PinnedTweetID: "1271186240323432452",
						},
						Tweet: &TweetObj{
							ID:   "1271186240323432452",
							Text: "We’re disclosing new state-linked information operations to our public archive — the only one of its kind in the industry. Originating from the People’s Republic of China (PRC), Russia, and Turkey, all associated accounts and content have been removed. https://t.co/obRqr96iYm",
						},
					},
					"3260518932": UserLookup{
						User: UserObj{
							ID:       "3260518932",
							Name:     "Twitter Moments",
							UserName: "TwitterMoments",
						},
					},
					"373471064": UserLookup{
						User: UserObj{
							ID:            "373471064",
							Name:          "Twitter Music",
							UserName:      "TwitterMusic",
							PinnedTweetID: "1293216056274759680",
						},
						Tweet: &TweetObj{
							ID:   "1293216056274759680",
							Text: "say howdy to your new yeehaw king @orvillepeck—our #ArtistToFollow this month 🤠 https://t.co/3pk9fYcPHb",
						},
					},
					"791978718": UserLookup{
						User: UserObj{
							ID:       "791978718",
							Name:     "Twitter Official Partner",
							UserName: "OfficialPartner",
						},
					},
					"17874544": UserLookup{
						User: UserObj{
							ID:            "17874544",
							Name:          "Twitter Support",
							UserName:      "TwitterSupport",
							PinnedTweetID: "1289000334497439744",
						},
						Tweet: &TweetObj{
							ID:   "1289000334497439744",
							Text: "We’ve significantly limited access to our internal tools and systems. Until we can safely resume normal operations, our response times to some support needs and reports will be slower. Thank you for your patience as we work through this.",
							ContextAnnotations: []TweetContextAnnotationObj{
								{
									Domain: TweetContextObj{
										ID:          "46",
										Name:        "Brand Category",
										Description: "Categories within Brand Verticals that narrow down the scope of Brands",
									},
									Entity: TweetContextObj{
										ID:   "781974596752842752",
										Name: "Services",
									},
								},
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
						},
					},
					"234489024": UserLookup{
						User: UserObj{
							ID:            "234489024",
							Name:          "Twitter Comms",
							UserName:      "TwitterComms",
							PinnedTweetID: "1283543147444711424",
						},
						Tweet: &TweetObj{
							ID:   "1283543147444711424",
							Text: "Follow @TwitterSupport for the latest on the security incident ⬇️ https://t.co/7FKKksJqxV",
							ContextAnnotations: []TweetContextAnnotationObj{
								{
									Domain: TweetContextObj{
										ID:          "46",
										Name:        "Brand Category",
										Description: "Categories within Brand Verticals that narrow down the scope of Brands",
									},
									Entity: TweetContextObj{
										ID:   "781974596752842752",
										Name: "Services",
									},
								},
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
						},
					},
					"1526228120": UserLookup{
						User: UserObj{
							ID:       "1526228120",
							Name:     "Twitter Data",
							UserName: "TwitterData",
						},
					},
				},
				Meta: &UserFollowMeta{
					ResultCount: 10,
					NextToken:   "DFEDBNRFT3MHCZZZ",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := u.LookupFollowing(context.Background(), tt.args.id, tt.args.followOpts)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.LookupFollowing() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User.LookupFollowing() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_LookupFollowers(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		id         string
		followOpts UserFollowOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UserFollowLookup
		wantErr bool
	}{
		{
			name: "Success-Basic",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), "2/users/2244994945/followers?max_results=10") == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetLookupEndpoint)
					}
					body := `{
						"data": [
						  {
							"id": "6253282",
							"name": "Twitter API",
							"username": "TwitterAPI"
						  },
						  {
							"id": "2244994945",
							"name": "Twitter Dev",
							"username": "TwitterDev"
						  },
						  {
							"id": "783214",
							"name": "Twitter",
							"username": "Twitter"
						  },
						  {
							"id": "95731075",
							"name": "Twitter Safety",
							"username": "TwitterSafety"
						  },
						  {
							"id": "3260518932",
							"name": "Twitter Moments",
							"username": "TwitterMoments"
						  },
						  {
							"id": "373471064",
							"name": "Twitter Music",
							"username": "TwitterMusic"
						  },
						  {
							"id": "791978718",
							"name": "Twitter Official Partner",
							"username": "OfficialPartner"
						  },
						  {
							"id": "17874544",
							"name": "Twitter Support",
							"username": "TwitterSupport"
						  },
						  {
							"id": "234489024",
							"name": "Twitter Comms",
							"username": "TwitterComms"
						  },
						  {
							"id": "1526228120",
							"name": "Twitter Data",
							"username": "TwitterData"
						  }
						],
						"meta": {
						  "result_count": 10,
						  "next_token": "DFEDBNRFT3MHCZZZ"
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				id: "2244994945",
				followOpts: UserFollowOptions{
					MaxResults: 10,
				},
			},
			want: &UserFollowLookup{
				Lookups: UserLookups{
					"6253282": UserLookup{
						User: UserObj{
							ID:       "6253282",
							Name:     "Twitter API",
							UserName: "TwitterAPI",
						},
					},
					"2244994945": UserLookup{
						User: UserObj{
							ID:       "2244994945",
							Name:     "Twitter Dev",
							UserName: "TwitterDev",
						},
					},
					"783214": UserLookup{
						User: UserObj{
							ID:       "783214",
							Name:     "Twitter",
							UserName: "Twitter",
						},
					},
					"95731075": UserLookup{
						User: UserObj{
							ID:       "95731075",
							Name:     "Twitter Safety",
							UserName: "TwitterSafety",
						},
					},
					"3260518932": UserLookup{
						User: UserObj{
							ID:       "3260518932",
							Name:     "Twitter Moments",
							UserName: "TwitterMoments",
						},
					},
					"373471064": UserLookup{
						User: UserObj{
							ID:       "373471064",
							Name:     "Twitter Music",
							UserName: "TwitterMusic",
						},
					},
					"791978718": UserLookup{
						User: UserObj{
							ID:       "791978718",
							Name:     "Twitter Official Partner",
							UserName: "OfficialPartner",
						},
					},
					"17874544": UserLookup{
						User: UserObj{
							ID:       "17874544",
							Name:     "Twitter Support",
							UserName: "TwitterSupport",
						},
					},
					"234489024": UserLookup{
						User: UserObj{
							ID:       "234489024",
							Name:     "Twitter Comms",
							UserName: "TwitterComms",
						},
					},
					"1526228120": UserLookup{
						User: UserObj{
							ID:       "1526228120",
							Name:     "Twitter Data",
							UserName: "TwitterData",
						},
					},
				},
				Meta: &UserFollowMeta{
					ResultCount: 10,
					NextToken:   "DFEDBNRFT3MHCZZZ",
				},
			},
			wantErr: false,
		},
		{
			name: "Success-Optional",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), "2/users/2244994945/followers?") == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetLookupEndpoint)
					}
					body := `{
						"data": [
						  {
							"pinned_tweet_id": "1293595870563381249",
							"id": "6253282",
							"username": "TwitterAPI",
							"name": "Twitter API"
						  },
						  {
							"pinned_tweet_id": "1293593516040269825",
							"id": "2244994945",
							"username": "TwitterDev",
							"name": "Twitter Dev"
						  },
						  {
							"id": "783214",
							"username": "Twitter",
							"name": "Twitter"
						  },
						  {
							"pinned_tweet_id": "1271186240323432452",
							"id": "95731075",
							"username": "TwitterSafety",
							"name": "Twitter Safety"
						  },
						  {
							"id": "3260518932",
							"username": "TwitterMoments",
							"name": "Twitter Moments"
						  },
						  {
							"pinned_tweet_id": "1293216056274759680",
							"id": "373471064",
							"username": "TwitterMusic",
							"name": "Twitter Music"
						  },
						  {
							"id": "791978718",
							"username": "OfficialPartner",
							"name": "Twitter Official Partner"
						  },
						  {
							"pinned_tweet_id": "1289000334497439744",
							"id": "17874544",
							"username": "TwitterSupport",
							"name": "Twitter Support"
						  },
						  {
							"pinned_tweet_id": "1283543147444711424",
							"id": "234489024",
							"username": "TwitterComms",
							"name": "Twitter Comms"
						  },
						  {
							"id": "1526228120",
							"username": "TwitterData",
							"name": "Twitter Data"
						  }
						],
						"includes": {
						  "tweets": [
							{
							  "context_annotations": [
								{
								  "domain": {
									"id": "46",
									"name": "Brand Category",
									"description": "Categories within Brand Verticals that narrow down the scope of Brands"
								  },
								  "entity": {
									"id": "781974596752842752",
									"name": "Services"
								  }
								},
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
								},
								{
								  "domain": {
									"id": "65",
									"name": "Interests and Hobbies Vertical",
									"description": "Top level interests and hobbies groupings, like Food or Travel"
								  },
								  "entity": {
									"id": "848920371311001600",
									"name": "Technology",
									"description": "Technology and computing"
								  }
								},
								{
								  "domain": {
									"id": "66",
									"name": "Interests and Hobbies Category",
									"description": "A grouping of interests and hobbies entities, like Novelty Food or Destinations"
								  },
								  "entity": {
									"id": "848921413196984320",
									"name": "Computer programming",
									"description": "Computer programming"
								  }
								},
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
							  "id": "1293595870563381249",
							  "text": "Twitter API v2: Early Access released\n\nToday we announced Early Access to the first endpoints of the new Twitter API!\n\n#TwitterAPI #EarlyAccess #VersionBump https://t.co/g7v3aeIbtQ"
							},
							{
							  "context_annotations": [
								{
								  "domain": {
									"id": "46",
									"name": "Brand Category",
									"description": "Categories within Brand Verticals that narrow down the scope of Brands"
								  },
								  "entity": {
									"id": "781974596752842752",
									"name": "Services"
								  }
								},
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
								},
								{
								  "domain": {
									"id": "65",
									"name": "Interests and Hobbies Vertical",
									"description": "Top level interests and hobbies groupings, like Food or Travel"
								  },
								  "entity": {
									"id": "848920371311001600",
									"name": "Technology",
									"description": "Technology and computing"
								  }
								},
								{
								  "domain": {
									"id": "66",
									"name": "Interests and Hobbies Category",
									"description": "A grouping of interests and hobbies entities, like Novelty Food or Destinations"
								  },
								  "entity": {
									"id": "848921413196984320",
									"name": "Computer programming",
									"description": "Computer programming"
								  }
								}
							  ],
							  "id": "1293593516040269825",
							  "text": "It’s finally here! 🥁 Say hello to the new #TwitterAPI.\n\nWe’re rebuilding the Twitter API v2 from the ground up to better serve our developer community. And today’s launch is only the beginning.\n\nhttps://t.co/32VrwpGaJw https://t.co/KaFSbjWUA8"
							},
							{
							  "id": "1271186240323432452",
							  "text": "We’re disclosing new state-linked information operations to our public archive — the only one of its kind in the industry. Originating from the People’s Republic of China (PRC), Russia, and Turkey, all associated accounts and content have been removed. https://t.co/obRqr96iYm"
							},
							{
							  "id": "1293216056274759680",
							  "text": "say howdy to your new yeehaw king @orvillepeck—our #ArtistToFollow this month 🤠 https://t.co/3pk9fYcPHb"
							},
							{
							  "context_annotations": [
								{
								  "domain": {
									"id": "46",
									"name": "Brand Category",
									"description": "Categories within Brand Verticals that narrow down the scope of Brands"
								  },
								  "entity": {
									"id": "781974596752842752",
									"name": "Services"
								  }
								},
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
							  "id": "1289000334497439744",
							  "text": "We’ve significantly limited access to our internal tools and systems. Until we can safely resume normal operations, our response times to some support needs and reports will be slower. Thank you for your patience as we work through this."
							},
							{
							  "context_annotations": [
								{
								  "domain": {
									"id": "46",
									"name": "Brand Category",
									"description": "Categories within Brand Verticals that narrow down the scope of Brands"
								  },
								  "entity": {
									"id": "781974596752842752",
									"name": "Services"
								  }
								},
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
							  "id": "1283543147444711424",
							  "text": "Follow @TwitterSupport for the latest on the security incident ⬇️ https://t.co/7FKKksJqxV"
							}
						  ]
						},
						"meta": {
							"result_count": 10,
							"next_token": "DFEDBNRFT3MHCZZZ"
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
				id: "2244994945",
				followOpts: UserFollowOptions{
					MaxResults:  10,
					Expansions:  []Expansion{ExpansionPinnedTweetID},
					TweetFields: []TweetField{TweetFieldID, TweetFieldContextAnnotations},
				},
			},
			want: &UserFollowLookup{
				Lookups: UserLookups{
					"6253282": UserLookup{
						User: UserObj{
							ID:            "6253282",
							Name:          "Twitter API",
							UserName:      "TwitterAPI",
							PinnedTweetID: "1293595870563381249",
						},
						Tweet: &TweetObj{
							ID:   "1293595870563381249",
							Text: "Twitter API v2: Early Access released\n\nToday we announced Early Access to the first endpoints of the new Twitter API!\n\n#TwitterAPI #EarlyAccess #VersionBump https://t.co/g7v3aeIbtQ",
							ContextAnnotations: []TweetContextAnnotationObj{
								{
									Domain: TweetContextObj{
										ID:          "46",
										Name:        "Brand Category",
										Description: "Categories within Brand Verticals that narrow down the scope of Brands",
									},
									Entity: TweetContextObj{
										ID:   "781974596752842752",
										Name: "Services",
									},
								},
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
								{
									Domain: TweetContextObj{
										ID:          "65",
										Name:        "Interests and Hobbies Vertical",
										Description: "Top level interests and hobbies groupings, like Food or Travel",
									},
									Entity: TweetContextObj{
										ID:          "848920371311001600",
										Name:        "Technology",
										Description: "Technology and computing",
									},
								},
								{
									Domain: TweetContextObj{
										ID:          "66",
										Name:        "Interests and Hobbies Category",
										Description: "A grouping of interests and hobbies entities, like Novelty Food or Destinations",
									},
									Entity: TweetContextObj{
										ID:          "848921413196984320",
										Name:        "Computer programming",
										Description: "Computer programming",
									},
								},
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
						},
					},
					"2244994945": UserLookup{
						User: UserObj{
							ID:            "2244994945",
							Name:          "Twitter Dev",
							UserName:      "TwitterDev",
							PinnedTweetID: "1293593516040269825",
						},
						Tweet: &TweetObj{
							ID:   "1293593516040269825",
							Text: "It’s finally here! 🥁 Say hello to the new #TwitterAPI.\n\nWe’re rebuilding the Twitter API v2 from the ground up to better serve our developer community. And today’s launch is only the beginning.\n\nhttps://t.co/32VrwpGaJw https://t.co/KaFSbjWUA8",
							ContextAnnotations: []TweetContextAnnotationObj{
								{
									Domain: TweetContextObj{
										ID:          "46",
										Name:        "Brand Category",
										Description: "Categories within Brand Verticals that narrow down the scope of Brands",
									},
									Entity: TweetContextObj{
										ID:   "781974596752842752",
										Name: "Services",
									},
								},
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
								{
									Domain: TweetContextObj{
										ID:          "65",
										Name:        "Interests and Hobbies Vertical",
										Description: "Top level interests and hobbies groupings, like Food or Travel",
									},
									Entity: TweetContextObj{
										ID:          "848920371311001600",
										Name:        "Technology",
										Description: "Technology and computing",
									},
								},
								{
									Domain: TweetContextObj{
										ID:          "66",
										Name:        "Interests and Hobbies Category",
										Description: "A grouping of interests and hobbies entities, like Novelty Food or Destinations",
									},
									Entity: TweetContextObj{
										ID:          "848921413196984320",
										Name:        "Computer programming",
										Description: "Computer programming",
									},
								},
							},
						},
					},
					"783214": UserLookup{
						User: UserObj{
							ID:       "783214",
							Name:     "Twitter",
							UserName: "Twitter",
						},
					},
					"95731075": UserLookup{
						User: UserObj{
							ID:            "95731075",
							Name:          "Twitter Safety",
							UserName:      "TwitterSafety",
							PinnedTweetID: "1271186240323432452",
						},
						Tweet: &TweetObj{
							ID:   "1271186240323432452",
							Text: "We’re disclosing new state-linked information operations to our public archive — the only one of its kind in the industry. Originating from the People’s Republic of China (PRC), Russia, and Turkey, all associated accounts and content have been removed. https://t.co/obRqr96iYm",
						},
					},
					"3260518932": UserLookup{
						User: UserObj{
							ID:       "3260518932",
							Name:     "Twitter Moments",
							UserName: "TwitterMoments",
						},
					},
					"373471064": UserLookup{
						User: UserObj{
							ID:            "373471064",
							Name:          "Twitter Music",
							UserName:      "TwitterMusic",
							PinnedTweetID: "1293216056274759680",
						},
						Tweet: &TweetObj{
							ID:   "1293216056274759680",
							Text: "say howdy to your new yeehaw king @orvillepeck—our #ArtistToFollow this month 🤠 https://t.co/3pk9fYcPHb",
						},
					},
					"791978718": UserLookup{
						User: UserObj{
							ID:       "791978718",
							Name:     "Twitter Official Partner",
							UserName: "OfficialPartner",
						},
					},
					"17874544": UserLookup{
						User: UserObj{
							ID:            "17874544",
							Name:          "Twitter Support",
							UserName:      "TwitterSupport",
							PinnedTweetID: "1289000334497439744",
						},
						Tweet: &TweetObj{
							ID:   "1289000334497439744",
							Text: "We’ve significantly limited access to our internal tools and systems. Until we can safely resume normal operations, our response times to some support needs and reports will be slower. Thank you for your patience as we work through this.",
							ContextAnnotations: []TweetContextAnnotationObj{
								{
									Domain: TweetContextObj{
										ID:          "46",
										Name:        "Brand Category",
										Description: "Categories within Brand Verticals that narrow down the scope of Brands",
									},
									Entity: TweetContextObj{
										ID:   "781974596752842752",
										Name: "Services",
									},
								},
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
						},
					},
					"234489024": UserLookup{
						User: UserObj{
							ID:            "234489024",
							Name:          "Twitter Comms",
							UserName:      "TwitterComms",
							PinnedTweetID: "1283543147444711424",
						},
						Tweet: &TweetObj{
							ID:   "1283543147444711424",
							Text: "Follow @TwitterSupport for the latest on the security incident ⬇️ https://t.co/7FKKksJqxV",
							ContextAnnotations: []TweetContextAnnotationObj{
								{
									Domain: TweetContextObj{
										ID:          "46",
										Name:        "Brand Category",
										Description: "Categories within Brand Verticals that narrow down the scope of Brands",
									},
									Entity: TweetContextObj{
										ID:   "781974596752842752",
										Name: "Services",
									},
								},
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
						},
					},
					"1526228120": UserLookup{
						User: UserObj{
							ID:       "1526228120",
							Name:     "Twitter Data",
							UserName: "TwitterData",
						},
					},
				},
				Meta: &UserFollowMeta{
					ResultCount: 10,
					NextToken:   "DFEDBNRFT3MHCZZZ",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := u.LookupFollowers(context.Background(), tt.args.id, tt.args.followOpts)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.LookupFollowers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User.LookupFollowers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_Tweets(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		id        string
		tweetOpts UserTimelineOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UserTimeline
		wantErr bool
	}{
		{
			name: "Success-Default",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), "2/users/2244994945/tweets") == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetLookupEndpoint)
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
						  },
						  {
							"id": "1337498609819021312",
							"text": "Thanks to everyone who tuned in today to make music with the #TwitterAPI!\n\nNext week on Twitch - @iamdaniele and @jessicagarson will show you how to integrate the #TwitterAPI and Google Sheets 📈. Tuesday, Dec 15th at 2pm ET. \n\nhttps://t.co/SQziic6eyp"
						  }
						],
						"meta": {
						  "oldest_id": "1334564488884862976",
						  "newest_id": "1338971066773905408",
						  "result_count": 10,
						  "next_token": "7140dibdnow9c7btw3w29grvxfcgvpb9n9coehpk7xz5i"
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				id: "2244994945",
			},
			want: &UserTimeline{
				Tweets: []TweetObj{
					{
						ID:   "1338971066773905408",
						Text: "💡 Using Twitter data for academic research? Join our next livestream this Friday @ 9am PT on https://t.co/GrtBOXh5Y1!\n \n@SuhemParack will show how to get started with recent search &amp; filtered stream endpoints on the #TwitterAPI v2, the new Tweet payload, annotations, &amp; more. https://t.co/IraD2Z7wEg",
					},
					{
						ID:   "1338923691497959425",
						Text: "📈 Live now with @jessicagarson and @i_am_daniele! https://t.co/Y1AFzsTTxb",
					},
					{
						ID:   "1337498609819021312",
						Text: "Thanks to everyone who tuned in today to make music with the #TwitterAPI!\n\nNext week on Twitch - @iamdaniele and @jessicagarson will show you how to integrate the #TwitterAPI and Google Sheets 📈. Tuesday, Dec 15th at 2pm ET. \n\nhttps://t.co/SQziic6eyp",
					},
				},
				Meta: UserTimelineMeta{
					OldestID:    "1334564488884862976",
					NewestID:    "1338971066773905408",
					ResultCount: 10,
					NextToken:   "7140dibdnow9c7btw3w29grvxfcgvpb9n9coehpk7xz5i",
				},
			},
		},
		{
			name: "Success-Optional",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), "2/users/2244994945/tweets?") == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetLookupEndpoint)
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
								  "id": "46",
								  "name": "Brand Category",
								  "description": "Categories within Brand Verticals that narrow down the scope of Brands"
								},
								"entity": {
								  "id": "781974596752842752",
								  "name": "Services"
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
						  "oldest_id": "1337122535188652033",
						  "newest_id": "1338971066773905408",
						  "result_count": 5,
						  "next_token": "7140dibdnow9c7btw3w29n4v1mtag9kegr0gr7y26pnw3"
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				id: "2244994945",
				tweetOpts: UserTimelineOpts{
					MaxResults:  5,
					TweetFields: []TweetField{TweetFieldCreatedAt, TweetFieldAuthorID, TweetFieldConversationID, TweetFieldPublicMetrics, TweetFieldContextAnnotations},
					UserFields:  []UserField{UserFieldName},
					Expansions:  []Expansion{ExpansionAuthorID},
				},
			},
			want: &UserTimeline{
				Tweets: []TweetObj{
					{
						AuthorID:       "2244994945",
						ConversationID: "1338971066773905408",
						ID:             "1338971066773905408",
						Text:           "💡 Using Twitter data for academic research? Join our next livestream this Friday @ 9am PT on https://t.co/GrtBOXh5Y1!\n \n@SuhemParack will show how to get started with recent search &amp; filtered stream endpoints on the #TwitterAPI v2, the new Tweet payload, annotations, &amp; more. https://t.co/IraD2Z7wEg",
						ContextAnnotations: []TweetContextAnnotationObj{
							{
								Domain: TweetContextObj{
									ID:          "46",
									Name:        "Brand Category",
									Description: "Categories within Brand Verticals that narrow down the scope of Brands",
								},
								Entity: TweetContextObj{
									ID:   "781974596752842752",
									Name: "Services",
								},
							},
						},
						PublicMetrics: TweetMetricsObj{
							Retweets: 10,
							Replies:  1,
							Likes:    41,
							Quotes:   4,
						},
						CreatedAt: "2020-12-15T22:15:53.000Z",
					},
				},
				Includes: &UserTimelineIncludes{
					Users: []UserObj{
						{
							ID:       "2244994945",
							Name:     "Twitter Dev",
							UserName: "TwitterDev",
						},
					},
				},
				Meta: UserTimelineMeta{
					OldestID:    "1337122535188652033",
					NewestID:    "1338971066773905408",
					ResultCount: 5,
					NextToken:   "7140dibdnow9c7btw3w29n4v1mtag9kegr0gr7y26pnw3",
				},
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
			got, err := u.Tweets(context.Background(), tt.args.id, tt.args.tweetOpts)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Tweets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User.Tweets() = %v, want %v", got, tt.want)
			}
		})
	}
}

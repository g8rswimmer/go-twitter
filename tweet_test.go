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

func TestTweet_Lookup(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		ids        []string
		parameters TweetFieldOptions
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           TweetLookups
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
					if strings.Contains(req.URL.String(), tweetLookupEndpoint) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetLookupEndpoint)
					}
					body := `{
						"data": {
						  "author_id": "2244994945",
						  "created_at": "2018-11-26T16:37:10.000Z",
						  "text": "Just getting started with Twitter APIs? Find out what you need in order to build an app. Watch this video! https://t.co/Hg8nkfoizN",
						  "id": "1067094924124872705"
						},
						"includes": {
						  "users": [
							{
							  "verified": true,
							  "username": "TwitterDev",
							  "id": "2244994945",
							  "name": "Twitter Dev"
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
				ids: []string{"1067094924124872705"},
				parameters: TweetFieldOptions{
					UserFields: []UserField{UserFieldVerified, UserFieldUserName, UserFieldID, UserFieldName},
				},
			},
			want: TweetLookups{
				"1067094924124872705": TweetLookup{
					Tweet: TweetObj{
						AuthorID:  "2244994945",
						CreatedAt: "2018-11-26T16:37:10.000Z",
						Text:      "Just getting started with Twitter APIs? Find out what you need in order to build an app. Watch this video! https://t.co/Hg8nkfoizN",
						ID:        "1067094924124872705",
					},
					User: &UserObj{
						Verified: true,
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
					if strings.Contains(req.URL.String(), tweetLookupEndpoint) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetLookupEndpoint)
					}
					body := `{
						"data": [
						  {
							"id": "1261326399320715264",
							"text": "Tune in to the @MongoDB @Twitch stream featuring our very own @suhemparack to learn about Twitter Developer Labs - starting now! https://t.co/fAWpYi3o5O",
							"author_id": "2244994945",
							"created_at": "2020-05-15T16:03:42.000Z"
						  },
						  {
							"id": "1278347468690915330",
							"text": "Good news and bad news: \n\n2020 is half over",
							"author_id": "783214",
							"created_at": "2020-07-01T15:19:21.000Z"
						  }
						],
						"includes": {
						  "users": [
							{
							  "verified": true,
							  "name": "Twitter Dev",
							  "id": "2244994945",
							  "username": "TwitterDev"
							},
							{
							  "verified": true,
							  "name": "Twitter",
							  "id": "783214",
							  "username": "Twitter"
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
				ids: []string{"1261326399320715264", "1278347468690915330"},
				parameters: TweetFieldOptions{
					UserFields: []UserField{UserFieldVerified, UserFieldUserName, UserFieldID, UserFieldName},
				},
			},
			want: TweetLookups{
				"1261326399320715264": TweetLookup{
					Tweet: TweetObj{
						AuthorID:  "2244994945",
						CreatedAt: "2020-05-15T16:03:42.000Z",
						Text:      "Tune in to the @MongoDB @Twitch stream featuring our very own @suhemparack to learn about Twitter Developer Labs - starting now! https://t.co/fAWpYi3o5O",
						ID:        "1261326399320715264",
					},
					User: &UserObj{
						Verified: true,
						UserName: "TwitterDev",
						ID:       "2244994945",
						Name:     "Twitter Dev",
					},
				},
				"1278347468690915330": TweetLookup{
					Tweet: TweetObj{
						AuthorID:  "783214",
						CreatedAt: "2020-07-01T15:19:21.000Z",
						Text:      "Good news and bad news: \n\n2020 is half over",
						ID:        "1278347468690915330",
					},
					User: &UserObj{
						Verified: true,
						UserName: "Twitter",
						ID:       "783214",
						Name:     "Twitter",
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
					if strings.Contains(req.URL.String(), tweetLookupEndpoint) == false {
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
				parameters: TweetFieldOptions{
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
			tl := &Tweet{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := tl.Lookup(context.Background(), tt.args.ids, tt.args.parameters)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tweet.Lookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tweet.Lookup() = %v, want %v", got, tt.want)
			}
			var tweetErr *TweetErrorResponse
			if errors.As(err, &tweetErr) {
				if !reflect.DeepEqual(tweetErr, tt.wantTweetError) {
					t.Errorf("Tweet.Lookup() Error = %+v, want %+v", tweetErr, tt.wantTweetError)
				}
			}
		})
	}
}

func TestTweet_RecentSearch(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		query      string
		searchOpts TweetRecentSearchOptions
		fieldOpts  TweetFieldOptions
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           *TweetRecentSearch
		wantErr        bool
		wantTweetError *TweetErrorResponse
	}{
		{
			name: "success query",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), tweetLookupEndpoint) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetRecentSearchEndpoint)
					}
					body := `{
						"data": [
						  {
							"id": "1279990139888918528",
							"text": "Python now online for you !!\n\nWith the advent and acceptance of AI, Robotics, Python has become an inevitable factor in software development industry and most looked out skill both Nationally and Internationally. \n\nCoupon code: GVUP9\nCall: 9482303905/9482163905 https://t.co/ZFXCDJedAh"
						  },
						  {
							"id": "1279990133463429120",
							"text": "RT @McQubit: Building Neural Networks with Python Code and Math in Detail — II https://t.co/l6PKTTFGkv #machine_learning #programming #math…"
						  },
						  {
							"id": "1279990118355476480",
							"text": "RT @SunnyVaram: Top 10 Natural Language Processing Online Courses https://t.co/oAGqkHdS8H via @https://twitter.com/analyticsinme \n#DataScie…"
						  },
						  {
							"id": "1279990114584875009",
							"text": "RT @mohitnihalani7: LINK IN BIO......\n\n#programming #coding #programmer #developer #python #code #technology #coder #javascript #java #comp…"
						  },
						  {
							"id": "1279990108968665088",
							"text": "RT @mohitnihalani7: LINK IN BIO......\n\n#programming #coding #programmer #developer #python #code #technology #coder #javascript #java #comp…"
						  },
						  {
							"id": "1279990090828320769",
							"text": "RT @SunnyVaram: Top 10 Natural Language Processing Online Courses https://t.co/oAGqkHdS8H via @https://twitter.com/analyticsinme \n#DataScie…"
						  },
						  {
							"id": "1279990084398387201",
							"text": "RT @mohitnihalani7: LINK IN BIO......\n\n#programming #coding #programmer #developer #python #code #technology #coder #javascript #java #comp…"
						  },
						  {
							"id": "1279990076748038145",
							"text": "RT @gp_pulipaka: Best Machine Learning and Data Science #Books 2020. #BigData #Analytics #DataScience #IoT #IIoT #PyTorch #Python #RStats #…"
						  },
						  {
							"id": "1279990069105917952",
							"text": "RT @SunnyVaram: Top 10 Natural Language Processing Online Courses https://t.co/oAGqkHdS8H via @https://twitter.com/analyticsinme \n#DataScie…"
						  },
						  {
							"id": "1279990063888076800",
							"text": "RT @mohitnihalani7: LINK IN BIO......\n\n#programming #coding #programmer #developer #python #code #technology #coder #javascript #java #comp…"
						  }
						],
						"meta": {
						  "newest_id": "1279990139888918528",
						  "oldest_id": "1279990063888076800",
						  "result_count": 10,
						  "next_token": "b26v89c19zqg8o3fo7gghep0wmpt92c0wn0jiqwtc7tdp"
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				query: "python",
				fieldOpts: TweetFieldOptions{
					UserFields: []UserField{UserFieldVerified, UserFieldUserName, UserFieldID, UserFieldName},
				},
			},
			want: &TweetRecentSearch{
				LookUps: TweetLookups{
					"1279990139888918528": TweetLookup{
						Tweet: TweetObj{
							ID:   "1279990139888918528",
							Text: "Python now online for you !!\n\nWith the advent and acceptance of AI, Robotics, Python has become an inevitable factor in software development industry and most looked out skill both Nationally and Internationally. \n\nCoupon code: GVUP9\nCall: 9482303905/9482163905 https://t.co/ZFXCDJedAh",
						},
					},
					"1279990133463429120": TweetLookup{
						Tweet: TweetObj{
							ID:   "1279990133463429120",
							Text: "RT @McQubit: Building Neural Networks with Python Code and Math in Detail — II https://t.co/l6PKTTFGkv #machine_learning #programming #math…",
						},
					},
					"1279990118355476480": TweetLookup{
						Tweet: TweetObj{
							ID:   "1279990118355476480",
							Text: "RT @SunnyVaram: Top 10 Natural Language Processing Online Courses https://t.co/oAGqkHdS8H via @https://twitter.com/analyticsinme \n#DataScie…",
						},
					},
					"1279990114584875009": TweetLookup{
						Tweet: TweetObj{
							ID:   "1279990114584875009",
							Text: "RT @mohitnihalani7: LINK IN BIO......\n\n#programming #coding #programmer #developer #python #code #technology #coder #javascript #java #comp…",
						},
					},
					"1279990108968665088": TweetLookup{
						Tweet: TweetObj{
							ID:   "1279990108968665088",
							Text: "RT @mohitnihalani7: LINK IN BIO......\n\n#programming #coding #programmer #developer #python #code #technology #coder #javascript #java #comp…",
						},
					},
					"1279990090828320769": TweetLookup{
						Tweet: TweetObj{
							ID:   "1279990090828320769",
							Text: "RT @SunnyVaram: Top 10 Natural Language Processing Online Courses https://t.co/oAGqkHdS8H via @https://twitter.com/analyticsinme \n#DataScie…",
						},
					},
					"1279990084398387201": TweetLookup{
						Tweet: TweetObj{
							ID:   "1279990084398387201",
							Text: "RT @mohitnihalani7: LINK IN BIO......\n\n#programming #coding #programmer #developer #python #code #technology #coder #javascript #java #comp…",
						},
					},
					"1279990076748038145": TweetLookup{
						Tweet: TweetObj{
							ID:   "1279990076748038145",
							Text: "RT @gp_pulipaka: Best Machine Learning and Data Science #Books 2020. #BigData #Analytics #DataScience #IoT #IIoT #PyTorch #Python #RStats #…",
						},
					},
					"1279990069105917952": TweetLookup{
						Tweet: TweetObj{
							ID:   "1279990069105917952",
							Text: "RT @SunnyVaram: Top 10 Natural Language Processing Online Courses https://t.co/oAGqkHdS8H via @https://twitter.com/analyticsinme \n#DataScie…",
						},
					},
					"1279990063888076800": TweetLookup{
						Tweet: TweetObj{
							ID:   "1279990063888076800",
							Text: "RT @mohitnihalani7: LINK IN BIO......\n\n#programming #coding #programmer #developer #python #code #technology #coder #javascript #java #comp…",
						},
					},
				},
				Meta: TweetRecentSearchMeta{
					NewestID:    "1279990139888918528",
					OldestID:    "1279990063888076800",
					ResultCount: 10,
					NextToken:   "b26v89c19zqg8o3fo7gghep0wmpt92c0wn0jiqwtc7tdp",
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
					if strings.Contains(req.URL.String(), tweetLookupEndpoint) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetRecentSearchEndpoint)
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
				query: "python",
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
			tweet := &Tweet{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := tweet.RecentSearch(context.Background(), tt.args.query, tt.args.searchOpts, tt.args.fieldOpts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tweet.RecentSearch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tweet.RecentSearch() = %v, want %v", got, tt.want)
			}
			var tweetErr *TweetErrorResponse
			if errors.As(err, &tweetErr) {
				if !reflect.DeepEqual(tweetErr, tt.wantTweetError) {
					t.Errorf("Tweet.Lookup() Error = %+v, want %+v", tweetErr, tt.wantTweetError)
				}
			}
		})
	}
}

func TestTweet_UpdateSearchStreamRules(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		rules    TweetSearchStreamRule
		validate bool
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           *TweetSearchStreamRules
		wantErr        bool
		wantTweetError *TweetErrorResponse
	}{
		{
			name: "Add rules",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodPost {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodPost)
					}
					if strings.Contains(req.URL.String(), tweetFilteredStreamRulesEndpoint) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetFilteredStreamRulesEndpoint)
					}

					body := `{
						"data": [
							{
								"value": "meme",
								"tag": "funny things",
								"id": "1166895166390583299"
							},
							{
								"value": "cats has:media -grumpy",
								"tag": "happy cats with media",
								"id": "1166895166390583296"
							},
							{
								"value": "cat has:media",
								"tag": "cats with media",
								"id": "1166895166390583297"
							},
							{
								"value": "meme has:images",
								"id": "1166895166390583298"
							}
					
						],
						"meta": {
							"sent": "2019-08-29T02:07:42.205Z",
							"summary": {
								"created": 4,
								"not_created": 0
							}
						}
					}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				rules: TweetSearchStreamRule{
					Add: []*TweetSearchStreamAddRule{
						{
							Value: "cats has:media",
							Tag:   "cats with media",
						},
						{
							Value: "cats has:media -grumpy",
							Tag:   "happy cats with media",
						},
						{
							Value: "meme",
							Tag:   "funny things",
						},
						{
							Value: "meme has:images",
						},
					},
				},
			},
			want: &TweetSearchStreamRules{
				Data: []TweetSearchStreamRuleData{
					{
						Value: "meme",
						Tag:   "funny things",
						ID:    "1166895166390583299",
					},
					{
						Value: "cats has:media -grumpy",
						Tag:   "happy cats with media",
						ID:    "1166895166390583296",
					},
					{
						Value: "cat has:media",
						Tag:   "cats with media",
						ID:    "1166895166390583297",
					},
					{
						Value: "meme has:images",
						ID:    "1166895166390583298",
					},
				},
				Meta: TweetSearchStreamRuleMeta{
					Sent: "2019-08-29T02:07:42.205Z",
					Summary: TweetSearchStreamRuleSummary{
						Created:    4,
						NotCreated: 0,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Delete rules",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodPost {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodPost)
					}
					if strings.Contains(req.URL.String(), tweetFilteredStreamRulesEndpoint) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetFilteredStreamRulesEndpoint)
					}
					body := `{
						"meta": {
						  "sent": "2019-08-29T01:48:54.633Z",
						  "summary": {
							"deleted": 1,
							"not_deleted": 0
						  }
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				rules: TweetSearchStreamRule{
					Delete: &TweetSearchStreamDeleteRule{
						IDs: []string{"1165037377523306498"},
					},
				},
			},
			want: &TweetSearchStreamRules{
				Meta: TweetSearchStreamRuleMeta{
					Sent: "2019-08-29T01:48:54.633Z",
					Summary: TweetSearchStreamRuleSummary{
						Deleted:    1,
						NotDeleted: 0,
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
					if req.Method != http.MethodPost {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodPost)
					}
					if strings.Contains(req.URL.String(), tweetFilteredStreamRulesEndpoint) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetFilteredStreamRulesEndpoint)
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
				rules: TweetSearchStreamRule{
					Add: []*TweetSearchStreamAddRule{
						{
							Value: "cats has:media",
							Tag:   "cats with media",
						},
						{
							Value: "cats has:media -grumpy",
							Tag:   "happy cats with media",
						},
						{
							Value: "meme",
							Tag:   "funny things",
						},
						{
							Value: "meme has:images",
						},
					},
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
			tweet := &Tweet{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := tweet.ApplyFilteredStreamRules(context.Background(), tt.args.rules, tt.args.validate)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tweet.UpdateSearchStreamRules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tweet.UpdateSearchStreamRules() = %+v, want %+v", got, tt.want)
			}
			var tweetErr *TweetErrorResponse
			if errors.As(err, &tweetErr) {
				if !reflect.DeepEqual(tweetErr, tt.wantTweetError) {
					t.Errorf("Tweet.Lookup() Error = %+v, want %+v", tweetErr, tt.wantTweetError)
				}
			}
		})
	}
}

func TestTweet_SearchStreamRules(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		ids []string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           *TweetSearchStreamRules
		wantErr        bool
		wantTweetError *TweetErrorResponse
	}{
		{
			name: "Get Rules",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), tweetFilteredStreamRulesEndpoint) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetFilteredStreamRulesEndpoint)
					}
					body := `{
						"data": [
						  {
							"id": "1165037377523306497",
							"value": "dog has:images",
							"tag": "dog pictures"
						  },
						  {
							"id": "1165037377523306498",
							"value": "cat has:images -grumpy"
						  }
						],
						"meta": {
						  "sent": "2019-08-29T01:12:10.729Z"
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				ids: []string{"1234"},
			},
			want: &TweetSearchStreamRules{
				Data: []TweetSearchStreamRuleData{
					{
						ID:    "1165037377523306497",
						Value: "dog has:images",
						Tag:   "dog pictures",
					},
					{
						ID:    "1165037377523306498",
						Value: "cat has:images -grumpy",
					},
				},
				Meta: TweetSearchStreamRuleMeta{
					Sent: "2019-08-29T01:12:10.729Z",
				},
			},
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
					if strings.Contains(req.URL.String(), tweetFilteredStreamRulesEndpoint) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetFilteredStreamRulesEndpoint)
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
				ids: []string{"1234"},
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
			tweet := &Tweet{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := tweet.FilteredStreamRules(context.Background(), tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tweet.SearchStreamRules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tweet.SearchStreamRules() = %v, want %v", got, tt.want)
			}
			var tweetErr *TweetErrorResponse
			if errors.As(err, &tweetErr) {
				if !reflect.DeepEqual(tweetErr, tt.wantTweetError) {
					t.Errorf("Tweet.Lookup() Error = %+v, want %+v", tweetErr, tt.wantTweetError)
				}
			}
		})
	}
}

func TestTweet_SearchStream(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		fieldOpts TweetFieldOptions
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           TweetLookups
		wantErr        bool
		wantTweetError *TweetErrorResponse
	}{
		{
			name: "search",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), tweetFilteredStreamEndpoint) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetFilteredStreamEndpoint)
					}
					body := `{
						"data": {
						  "id": "1067094924124872705",
						  "text": "Just getting started with Twitter APIs? Find out what you need in order to build an app. Watch this video! https://t.co/Hg8nkfoizN"
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				fieldOpts: TweetFieldOptions{
					UserFields: []UserField{UserFieldVerified, UserFieldUserName, UserFieldID, UserFieldName},
				},
			},
			want: TweetLookups{
				"1067094924124872705": TweetLookup{
					Tweet: TweetObj{
						Text: "Just getting started with Twitter APIs? Find out what you need in order to build an app. Watch this video! https://t.co/Hg8nkfoizN",
						ID:   "1067094924124872705",
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
					if strings.Contains(req.URL.String(), tweetFilteredStreamEndpoint) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetFilteredStreamEndpoint)
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
			args:    args{},
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
			tweet := &Tweet{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := tweet.FilteredStream(context.Background(), tt.args.fieldOpts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tweet.SearchStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tweet.SearchStream() = %v, want %v", got, tt.want)
			}
			var tweetErr *TweetErrorResponse
			if errors.As(err, &tweetErr) {
				if !reflect.DeepEqual(tweetErr, tt.wantTweetError) {
					t.Errorf("Tweet.Lookup() Error = %+v, want %+v", tweetErr, tt.wantTweetError)
				}
			}
		})
	}
}

func TestTweet_SampledStream(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		fieldOpts TweetFieldOptions
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           TweetLookups
		wantErr        bool
		wantTweetError *TweetErrorResponse
	}{
		{
			name: "search",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), tweetSampledStreamEndpoint) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetSampledStreamEndpoint)
					}
					body := `{
						"data": {
						  "id": "1067094924124872705",
						  "text": "Just getting started with Twitter APIs? Find out what you need in order to build an app. Watch this video! https://t.co/Hg8nkfoizN"
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				fieldOpts: TweetFieldOptions{
					UserFields: []UserField{UserFieldVerified, UserFieldUserName, UserFieldID, UserFieldName},
				},
			},
			want: TweetLookups{
				"1067094924124872705": TweetLookup{
					Tweet: TweetObj{
						Text: "Just getting started with Twitter APIs? Find out what you need in order to build an app. Watch this video! https://t.co/Hg8nkfoizN",
						ID:   "1067094924124872705",
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
					if strings.Contains(req.URL.String(), tweetSampledStreamEndpoint) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetSampledStreamEndpoint)
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
			args:    args{},
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
			tweet := &Tweet{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := tweet.SampledStream(context.Background(), tt.args.fieldOpts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tweet.SampledStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tweet.SampledStream() = %v, want %v", got, tt.want)
			}
			var tweetErr *TweetErrorResponse
			if errors.As(err, &tweetErr) {
				if !reflect.DeepEqual(tweetErr, tt.wantTweetError) {
					t.Errorf("Tweet.Lookup() Error = %+v, want %+v", tweetErr, tt.wantTweetError)
				}
			}
		})
	}
}

func TestTweet_Hide(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		id     string
		hidden bool
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		wantTweetError *TweetErrorResponse
	}{
		{
			name: "hide",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodPut {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodPut)
					}
					if strings.Contains(req.URL.String(), "hidden") == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), "hidden")
					}
					body := `{"data":{"hidden":true}}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				id:     "122334433",
				hidden: true,
			},
			wantErr: false,
		},
		{
			name: "unhide",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodPut {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodPut)
					}
					if strings.Contains(req.URL.String(), "hidden") == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), "hidden")
					}
					body := `{"data":{"hidden":false}}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				id:     "122334433",
				hidden: false,
			},
			wantErr: false,
		},
		{
			name: "mis-match",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodPut {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodPut)
					}
					if strings.Contains(req.URL.String(), "hidden") == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), "hidden")
					}
					body := `{"data":{"hidden":false}}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				id:     "122334433",
				hidden: true,
			},
			wantErr: true,
		},
		{
			name: "tweet error",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodPut {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodPut)
					}
					if strings.Contains(req.URL.String(), "hidden") == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), "hidden")
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
				id:     "122334433",
				hidden: false,
			},
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
			tweet := &Tweet{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			err := tweet.HideReplies(context.Background(), tt.args.id, tt.args.hidden)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tweet.Hide() error = %v, wantErr %v", err, tt.wantErr)
			}
			var tweetErr *TweetErrorResponse
			if errors.As(err, &tweetErr) {
				if !reflect.DeepEqual(tweetErr, tt.wantTweetError) {
					t.Errorf("Tweet.Lookup() Error = %+v, want %+v", tweetErr, tt.wantTweetError)
				}
			}
		})
	}
}

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

func TestClient_UserTweetTimeline(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		userID string
		opts   UserTweetTimelineOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UserTweetTimelineResponse
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
					if strings.Contains(req.URL.String(), userTweetTimelineEndpoint.urlID("", "2244994945")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userTweetTimelineEndpoint)
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
						  "oldest_id": "1334564488884862976",
						  "newest_id": "1338971066773905408",
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
			want: &UserTweetTimelineResponse{
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
				Meta: &UserTimelineMeta{
					ResultCount: 2,
					OldestID:    "1334564488884862976",
					NewestID:    "1338971066773905408",
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
					if strings.Contains(req.URL.String(), userTweetTimelineEndpoint.urlID("", "2244994945")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userTweetTimelineEndpoint)
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
						  "oldest_id": "1337122535188652033",
						  "newest_id": "1338971066773905408",
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
				opts: UserTweetTimelineOpts{
					TweetFields: []TweetField{TweetFieldCreatedAt, TweetFieldAuthorID, TweetFieldConversationID, TweetFieldPublicMetrics, TweetFieldContextAnnotations},
					UserFields:  []UserField{UserFieldUserName},
					Expansions:  []Expansion{ExpansionAuthorID},
					MaxResults:  10,
				},
			},
			want: &UserTweetTimelineResponse{
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
				Meta: &UserTimelineMeta{
					ResultCount: 2,
					OldestID:    "1337122535188652033",
					NewestID:    "1338971066773905408",
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
			got, err := c.UserTweetTimeline(context.Background(), tt.args.userID, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UserTweetTimeline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.UserTweetTimeline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_UserTweetReverseChronologicalTimeline(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		userID string
		opts   UserTweetChronologicalReverseTimelineOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UserTweetChronologicalReverseTimelineResponse
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
					if strings.Contains(req.URL.String(), userTweetReverseChronologicalTimelineEndpoint.urlID("", "2244994945")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userTweetReverseChronologicalTimelineEndpoint)
					}
					body := `{
						"data": [
						  {
							"id": "1524796546306478083",
							"text": "Today marks the launch of Devs in the Details, a technical video series made for developers by developers building with the Twitter API.  🚀nnIn this premiere episode, @jessicagarson walks us through how she built @FactualCat #WelcomeToOurTechTalkn⬇️nnhttps://t.co/nGa8JTQVBJ"
						  },
						  {
							"id": "1522642323535847424",
							"text": "We’ve gone into more detail on each Insider in our forum post. nnJoin us in congratulating the new additions! 🥳nnhttps://t.co/0r5maYEjPJ"
						  }
						],
						"meta": {
						  "result_count": 5,
						  "newest_id": "1524796546306478083",
						  "oldest_id": "1522642323535847424",
						  "next_token": "7140dibdnow9c7btw421dyz6jism75z99gyxd8egarsc4"
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
			want: &UserTweetChronologicalReverseTimelineResponse{
				Raw: &TweetRaw{
					Tweets: []*TweetObj{
						{
							ID:   "1524796546306478083",
							Text: "Today marks the launch of Devs in the Details, a technical video series made for developers by developers building with the Twitter API.  🚀nnIn this premiere episode, @jessicagarson walks us through how she built @FactualCat #WelcomeToOurTechTalkn⬇️nnhttps://t.co/nGa8JTQVBJ",
						},
						{
							ID:   "1522642323535847424",
							Text: "We’ve gone into more detail on each Insider in our forum post. nnJoin us in congratulating the new additions! 🥳nnhttps://t.co/0r5maYEjPJ",
						},
					},
				},
				Meta: &UserChronologicalReverseTimelineMeta{
					ResultCount: 5,
					NewestID:    "1524796546306478083",
					OldestID:    "1522642323535847424",
					NextToken:   "7140dibdnow9c7btw421dyz6jism75z99gyxd8egarsc4",
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
			name: "success with options",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), userTweetReverseChronologicalTimelineEndpoint.urlID("", "2244994945")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userTweetReverseChronologicalTimelineEndpoint)
					}
					body := `{
						"data": [
						  {
							"created_at": "2022-05-12T17:00:00.000Z",
							"text": "Today marks the launch of Devs in the Details, a technical video series made for developers by developers building with the Twitter API.  🚀nnIn this premiere episode, @jessicagarson walks us through how she built @FactualCat #WelcomeToOurTechTalkn⬇️nnhttps://t.co/nGa8JTQVBJ",
							"author_id": "2244994945",
							"id": "1524796546306478083"
						  },
						  {
							"created_at": "2022-05-06T18:19:53.000Z",
							"text": "We’ve gone into more detail on each Insider in our forum post. nnJoin us in congratulating the new additions! 🥳nnhttps://t.co/0r5maYEjPJ",
							"author_id": "2244994945",
							"id": "1522642323535847424"
   						  }
						],
						"includes": {
							"users": [
							  {
								"created_at": "2013-12-14T04:35:55.000Z",
								"name": "Twitter Dev",
								"username": "TwitterDev",
								"id": "2244994945"
							  }
							]
						},
						"meta": {
						  "result_count": 5,
						  "newest_id": "1524796546306478083",
						  "oldest_id": "1522642323535847424",
						  "next_token": "7140dibdnow9c7btw421dyz6jism75z99gyxd8egarsc4"
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
				opts: UserTweetChronologicalReverseTimelineOpts{
					Expansions:  []Expansion{ExpansionAuthorID},
					TweetFields: []TweetField{TweetFieldCreatedAt},
					UserFields:  []UserField{UserFieldCreatedAt, UserFieldName},
					MaxResults:  5,
				},
			},
			want: &UserTweetChronologicalReverseTimelineResponse{
				Raw: &TweetRaw{
					Tweets: []*TweetObj{
						{
							CreatedAt: "2022-05-12T17:00:00.000Z",
							AuthorID:  "2244994945",
							ID:        "1524796546306478083",
							Text:      "Today marks the launch of Devs in the Details, a technical video series made for developers by developers building with the Twitter API.  🚀nnIn this premiere episode, @jessicagarson walks us through how she built @FactualCat #WelcomeToOurTechTalkn⬇️nnhttps://t.co/nGa8JTQVBJ",
						},
						{
							CreatedAt: "2022-05-06T18:19:53.000Z",
							AuthorID:  "2244994945",
							ID:        "1522642323535847424",
							Text:      "We’ve gone into more detail on each Insider in our forum post. nnJoin us in congratulating the new additions! 🥳nnhttps://t.co/0r5maYEjPJ",
						},
					},
					Includes: &TweetRawIncludes{
						Users: []*UserObj{
							{
								CreatedAt: "2013-12-14T04:35:55.000Z",
								Name:      "Twitter Dev",
								UserName:  "TwitterDev",
								ID:        "2244994945",
							},
						},
					},
				},
				Meta: &UserChronologicalReverseTimelineMeta{
					ResultCount: 5,
					NewestID:    "1524796546306478083",
					OldestID:    "1522642323535847424",
					NextToken:   "7140dibdnow9c7btw421dyz6jism75z99gyxd8egarsc4",
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
			got, err := c.UserTweetReverseChronologicalTimeline(context.Background(), tt.args.userID, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UserTweetReverseChronologicalTimeline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.UserTweetReverseChronologicalTimeline() = %v, want %v", got, tt.want)
			}
		})
	}
}

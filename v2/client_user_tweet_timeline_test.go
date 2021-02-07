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
						Body:       ioutil.NopCloser(strings.NewReader(body)),
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
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				userID: "2244994945",
				opts: UserTweetTimelineOpts{
					TweetFields: []TweetField{TweetFieldCreatedAt, TweetFieldAuthorID, TweetFieldConversationID, TweetFieldPublicMetrics, TweetFieldContextAnnotations},
					UserFields:  []UserField{UserFieldUserName},
					Expansions:  []Expansion{ExpansionAuthorID},
					MaxResults:  2,
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

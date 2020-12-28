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

func TestClient_TweetLookup(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		ids  []string
		opts TweetLookupOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TweetLookupResponse
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
					if strings.Contains(req.URL.String(), string(tweetLookupEndpoint)) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetLookupEndpoint)
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
				ids: []string{"1067094924124872705"},
			},
			want: &TweetLookupResponse{
				Raw: &TweetRaw{
					Tweets: []*TweetObj{
						{
							Text: "Just getting started with Twitter APIs? Find out what you need in order to build an app. Watch this video! https://t.co/Hg8nkfoizN",
							ID:   "1067094924124872705",
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
					if strings.Contains(req.URL.String(), string(tweetLookupEndpoint)) == false {
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
				opts: TweetLookupOpts{
					Expansions:  []Expansion{ExpansionAttachmentsMediaKeys},
					MediaFields: []MediaField{MediaFieldType, MediaFieldDurationMS},
				},
			},
			want: &TweetLookupResponse{
				Raw: &TweetRaw{
					Tweets: []*TweetObj{
						{
							Text:      "Just getting started with Twitter APIs? Find out what you need in order to build an app. Watch this video! https://t.co/Hg8nkfoizN",
							ID:        "1067094924124872705",
							AuthorID:  "2244994945",
							CreatedAt: "2018-11-26T16:37:10.000Z",
						},
					},
					Includes: &TweetRawIncludes{
						Users: []*UserObj{
							{
								ID:       "2244994945",
								Verified: true,
								UserName: "TwitterDev",
								Name:     "Twitter Dev",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Success - Multiple IDs Default",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), string(tweetLookupEndpoint)) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetLookupEndpoint)
					}
					body := `{
						"data": [
						  {
							"id": "1261326399320715264",
							"text": "Tune in to the @MongoDB @Twitch stream featuring our very own @suhemparack to learn about Twitter Developer Labs - starting now! https://t.co/fAWpYi3o5O"
						  },
						  {
							"id": "1278347468690915330",
							"text": "Good news and bad news: \n\n2020 is half over"
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
				ids: []string{"1261326399320715264", "1278347468690915330"},
			},
			want: &TweetLookupResponse{
				Raw: &TweetRaw{
					Tweets: []*TweetObj{
						{
							Text: "Tune in to the @MongoDB @Twitch stream featuring our very own @suhemparack to learn about Twitter Developer Labs - starting now! https://t.co/fAWpYi3o5O",
							ID:   "1261326399320715264",
						},
						{
							Text: "Good news and bad news: \n\n2020 is half over",
							ID:   "1278347468690915330",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Success - Multiple IDs Optional",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), string(tweetLookupEndpoint)) == false {
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
				opts: TweetLookupOpts{
					Expansions:  []Expansion{ExpansionAuthorID},
					TweetFields: []TweetField{TweetFieldCreatedAt},
					UserFields:  []UserField{UserFieldName, UserFieldVerified},
				},
			},
			want: &TweetLookupResponse{
				Raw: &TweetRaw{
					Tweets: []*TweetObj{
						{
							Text:      "Tune in to the @MongoDB @Twitch stream featuring our very own @suhemparack to learn about Twitter Developer Labs - starting now! https://t.co/fAWpYi3o5O",
							ID:        "1261326399320715264",
							AuthorID:  "2244994945",
							CreatedAt: "2020-05-15T16:03:42.000Z",
						},
						{
							Text:      "Good news and bad news: \n\n2020 is half over",
							ID:        "1278347468690915330",
							AuthorID:  "783214",
							CreatedAt: "2020-07-01T15:19:21.000Z",
						},
					},
					Includes: &TweetRawIncludes{
						Users: []*UserObj{
							{
								ID:       "2244994945",
								Verified: true,
								UserName: "TwitterDev",
								Name:     "Twitter Dev",
							},
							{
								ID:       "783214",
								Verified: true,
								UserName: "Twitter",
								Name:     "Twitter",
							},
						},
					},
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
					if strings.Contains(req.URL.String(), string(tweetLookupEndpoint)) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetLookupEndpoint)
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
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				ids: []string{"1067094924124872705"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Success - Partial Errors",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), string(tweetLookupEndpoint)) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetLookupEndpoint)
					}
					body := `{
						"data": [
							{
						  		"id": "1067094924124872705",
						  		"text": "Just getting started with Twitter APIs? Find out what you need in order to build an app. Watch this video! https://t.co/Hg8nkfoizN"
							}
						],
						"errors": [
							{
							  "detail": "Could not find tweet with ids: [1276230436478386177].",
							  "title": "Not Found Error",
							  "resource_type": "tweet",
							  "parameter": "ids",
							  "value": "1276230436478386177",
							  "type": "https://api.twitter.com/2/problems/resource-not-found"
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
				ids: []string{"1067094924124872705", "1276230436478386177"},
			},
			want: &TweetLookupResponse{
				Raw: &TweetRaw{
					Tweets: []*TweetObj{
						{
							Text: "Just getting started with Twitter APIs? Find out what you need in order to build an app. Watch this video! https://t.co/Hg8nkfoizN",
							ID:   "1067094924124872705",
						},
					},
					Errors: []*ErrorObj{
						{
							Detail:       "Could not find tweet with ids: [1276230436478386177].",
							Title:        "Not Found Error",
							ResourceType: "tweet",
							Parameter:    "ids",
							Value:        "1276230436478386177",
							Type:         "https://api.twitter.com/2/problems/resource-not-found",
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
			got, err := c.TweetLookup(context.Background(), tt.args.ids, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.TweetLookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.TweetLookup() = %v, want %v", got, tt.want)
			}
		})
	}
}

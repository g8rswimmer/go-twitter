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

func TestClient_UserFollowersLookup(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		id   string
		opts UserFollowersLookupOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UserFollowersLookupResponse
		wantErr bool
	}{
		{
			name: "Success - Basic",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), userFollowersEndpoint.urlID("", "2244994945")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userFollowingEndpoint)
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
						  }
						],
						"meta": {
						  "result_count": 2,
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
			},
			want: &UserFollowersLookupResponse{
				Raw: &UserRaw{
					Users: []*UserObj{
						{
							ID:       "6253282",
							Name:     "Twitter API",
							UserName: "TwitterAPI",
						},
						{
							ID:       "2244994945",
							Name:     "Twitter Dev",
							UserName: "TwitterDev",
						},
					},
				},
				Meta: &UserFollowershMeta{
					ResultCount: 2,
					NextToken:   "DFEDBNRFT3MHCZZZ",
				},
			},
			wantErr: false,
		},
		{
			name: "Success - Options",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), userFollowersEndpoint.urlID("", "2244994945")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userFollowingEndpoint)
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
						  }
						],
						"includes": {
						  "tweets": [
							{
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
							  "id": "1293595870563381249",
							  "text": "Twitter API v2: Early Access released\n\nToday we announced Early Access to the first endpoints of the new Twitter API!\n\n#TwitterAPI #EarlyAccess #VersionBump https://t.co/g7v3aeIbtQ"
							},
							{
							  "id": "1293593516040269825",
							  "text": "We’re disclosing new state-linked information operations to our public archive — the only one of its kind in the industry. Originating from the People’s Republic of China (PRC), Russia, and Turkey, all associated accounts and content have been removed. https://t.co/obRqr96iYm"
							}
						  ]
						},
						"meta": {
							"result_count": 2,
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
				opts: UserFollowersLookupOpts{
					Expansions:  []Expansion{ExpansionPinnedTweetID},
					TweetFields: []TweetField{TweetFieldContextAnnotations},
				},
			},
			want: &UserFollowersLookupResponse{
				Raw: &UserRaw{
					Users: []*UserObj{
						{
							ID:            "6253282",
							Name:          "Twitter API",
							UserName:      "TwitterAPI",
							PinnedTweetID: "1293595870563381249",
						},
						{
							ID:            "2244994945",
							Name:          "Twitter Dev",
							UserName:      "TwitterDev",
							PinnedTweetID: "1293593516040269825",
						},
					},
					Includes: &UserRawIncludes{
						Tweets: []*TweetObj{
							{
								ID:   "1293595870563381249",
								Text: "Twitter API v2: Early Access released\n\nToday we announced Early Access to the first endpoints of the new Twitter API!\n\n#TwitterAPI #EarlyAccess #VersionBump https://t.co/g7v3aeIbtQ",
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
							},
							{
								ID:   "1293593516040269825",
								Text: "We’re disclosing new state-linked information operations to our public archive — the only one of its kind in the industry. Originating from the People’s Republic of China (PRC), Russia, and Turkey, all associated accounts and content have been removed. https://t.co/obRqr96iYm",
							},
						},
					},
				},
				Meta: &UserFollowershMeta{
					ResultCount: 2,
					NextToken:   "DFEDBNRFT3MHCZZZ",
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
			got, err := c.UserFollowersLookup(context.Background(), tt.args.id, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UserFollowersLookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.UserFollowersLookup() = %v, want %v", got, tt.want)
			}
		})
	}
}

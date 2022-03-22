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

func TestClient_SpacesLookup(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		ids  []string
		opts SpacesLookupOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *SpacesLookupResponse
		wantErr bool
	}{
		{
			name: "single success",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), string(spaceLookupEndpoint)+"/1DXxyRYNejbKM") == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), spaceLookupEndpoint)
					}
					body := `{
						"data": {
						  "id": "1DXxyRYNejbKM",
						  "state": "live"
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
				ids: []string{"1DXxyRYNejbKM"},
			},
			want: &SpacesLookupResponse{
				Raw: &SpacesRaw{
					Spaces: []*SpaceObj{
						{
							ID:    "1DXxyRYNejbKM",
							State: "live",
						},
					},
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
			name: "single success with option",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), string(spaceLookupEndpoint)+"/1DXxyRYNejbKM") == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), spaceLookupEndpoint)
					}
					body := `{
						"data": {
						  "host_ids": [
							"872212934402899973"
						  ],
						  "id": "1DXxyRYNejbKM",
						  "state": "live"
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				ids: []string{"1DXxyRYNejbKM"},
				opts: SpacesLookupOpts{
					SpaceFields: []SpaceField{SpaceFieldHostIDs},
				},
			},
			want: &SpacesLookupResponse{
				Raw: &SpacesRaw{
					Spaces: []*SpaceObj{
						{
							ID:    "1DXxyRYNejbKM",
							State: "live",
							HostIDs: []string{
								"872212934402899973",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "single success with option",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), string(spaceLookupEndpoint)) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), spaceLookupEndpoint)
					}
					body := `{
						"data": [
						  {
							"host_ids": [
							  "2244994945"
							],
							"id": "1DXxyRYNejbKM",
							"state": "live"
						  },
						  {
							"host_ids": [
							  "6253282"
							],
							"id": "1nAJELYEEPvGL",
							"state": "scheduled"
						  }
						]
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				ids: []string{"1DXxyRYNejbKM", "1nAJELYEEPvGL"},
				opts: SpacesLookupOpts{
					SpaceFields: []SpaceField{SpaceFieldHostIDs},
				},
			},
			want: &SpacesLookupResponse{
				Raw: &SpacesRaw{
					Spaces: []*SpaceObj{
						{
							ID:    "1DXxyRYNejbKM",
							State: "live",
							HostIDs: []string{
								"2244994945",
							},
						},
						{
							ID:    "1nAJELYEEPvGL",
							State: "scheduled",
							HostIDs: []string{
								"6253282",
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
			got, err := c.SpacesLookup(context.Background(), tt.args.ids, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SpacesLookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.SpacesLookup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_SpacesByCreatorLookup(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		userIDs []string
		opts    SpacesByCreatorLookupOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *SpacesByCreatorLookupResponse
		wantErr bool
	}{
		{
			name: "single success",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), string(spaceByCreatorLookupEndpoint)) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), spaceByCreatorLookupEndpoint)
					}
					body := `{
						"data": [
						  {
							"id": "1DXxyRYNejbKM",
							"state": "live"
						  }
						],
						"meta": {
						  "result_count": 1
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
				userIDs: []string{"2244994945"},
			},
			want: &SpacesByCreatorLookupResponse{
				Raw: &SpacesRaw{
					Spaces: []*SpaceObj{
						{
							ID:    "1DXxyRYNejbKM",
							State: "live",
						},
					},
				},
				Meta: &SpacesByCreatorMeta{
					ResultCount: 1,
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
			name: "single success with options",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), string(spaceByCreatorLookupEndpoint)) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), spaceByCreatorLookupEndpoint)
					}
					body := `{
						"data": [
						  {
							"host_ids": [
							  "2244994945"
							],
							"id": "1DXxyRYNejbKM",
							"state": "live"
						  },
						  {
							"host_ids": [
							  "6253282"
							],
							"id": "1nAJELYEEPvGL",
							"state": "scheduled"
						  }
						],
						"meta": {
						  "result_count": 2
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
				userIDs: []string{"2244994945", "6253282"},
				opts: SpacesByCreatorLookupOpts{
					SpaceFields: []SpaceField{SpaceFieldHostIDs},
				},
			},
			want: &SpacesByCreatorLookupResponse{
				Raw: &SpacesRaw{
					Spaces: []*SpaceObj{
						{
							ID:    "1DXxyRYNejbKM",
							State: "live",
							HostIDs: []string{
								"2244994945",
							},
						},
						{
							ID:    "1nAJELYEEPvGL",
							State: "scheduled",
							HostIDs: []string{
								"6253282",
							},
						},
					},
				},
				Meta: &SpacesByCreatorMeta{
					ResultCount: 2,
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
			got, err := c.SpacesByCreatorLookup(context.Background(), tt.args.userIDs, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SpacesByCreatorLookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.SpacesByCreatorLookup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_SpaceBuyersLookup(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		spaceID string
		opts    SpaceBuyersLookupOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *SpaceBuyersLookupResponse
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
					if strings.Contains(req.URL.String(), spaceBuyersLookupEndpoint.urlID("", "1DXxyRYNejbKM")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), spaceBuyersLookupEndpoint)
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
							  "text": "During these unprecedented times, what’s happening on Twitter can help the world better understand &amp; respond to the pandemic. nnWe're launching a free COVID-19 stream endpoint so qualified devs &amp; researchers can study the public conversation in real-time. https://t.co/BPqMcQzhId",
							  "id": "1255542774432063488"
							},
							{
							  "created_at": "2020-06-19T21:12:30.000Z",
							  "text": "📍 Minneapolisn🗣️ @FredTJoseph https://t.co/lNTOkyguG1",
							  "id": "1274087687469715457"
							}
						  ]
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
				spaceID: "1DXxyRYNejbKM",
				opts: SpaceBuyersLookupOpts{
					Expansions:  []Expansion{ExpansionPinnedTweetID},
					UserFields:  []UserField{UserFieldCreatedAt},
					TweetFields: []TweetField{TweetFieldCreatedAt},
				},
			},
			want: &SpaceBuyersLookupResponse{
				Raw: &UserRaw{
					Users: []*UserObj{
						{
							CreatedAt:     "2013-12-14T04:35:55.000Z",
							UserName:      "TwitterDev",
							PinnedTweetID: "1255542774432063488",
							ID:            "2244994945",
							Name:          "Twitter Dev",
						},
						{
							CreatedAt:     "2007-02-20T14:35:54.000Z",
							UserName:      "Twitter",
							PinnedTweetID: "1274087687469715457",
							ID:            "783214",
							Name:          "Twitter",
						},
					},
					Includes: &UserRawIncludes{
						Tweets: []*TweetObj{
							{
								CreatedAt: "2020-04-29T17:01:38.000Z",
								Text:      "During these unprecedented times, what’s happening on Twitter can help the world better understand &amp; respond to the pandemic. nnWe're launching a free COVID-19 stream endpoint so qualified devs &amp; researchers can study the public conversation in real-time. https://t.co/BPqMcQzhId",
								ID:        "1255542774432063488",
							},
							{
								CreatedAt: "2020-06-19T21:12:30.000Z",
								Text:      "📍 Minneapolisn🗣️ @FredTJoseph https://t.co/lNTOkyguG1",
								ID:        "1274087687469715457",
							},
						},
					},
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
			got, err := c.SpaceBuyersLookup(context.Background(), tt.args.spaceID, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SpaceBuyersLookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.SpaceBuyersLookup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_SpaceTweetsLookup(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		spaceID string
		opts    SpaceTweetsLookupOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *SpaceTweetsLookupResponse
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
					if strings.Contains(req.URL.String(), spaceTweetsLookupEndpoint.urlID("", "1DXxyRYNejbKM")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), spaceTweetsLookupEndpoint)
					}
					body := `{
						"data": [
						  {
							"id": "1389270063807598594",
							"author_id": "1065249714214457345",
							"text": "now, everyone with 600 or more followers can host a Space.nnbased on what we've learned, these accounts are likely to have a good experience hosting because of their existing audience. before bringing the ability to create a Space to everyone, we're focused on a few things. :thread:"
						  },
						  {
							"id": "1354143047324299264",
							"author_id": "783214",
							"text": "Academics are one of the biggest groups using the #TwitterAPI to research what's happening. Their work helps make the world (&amp; Twitter) a better place, and now more than ever, we must enable more of it. nIntroducing :drum_with_drumsticks: the Academic Research product track!nhttps://t.co/nOFiGewAV2"
						  },
						  {
							"id": "1293595870563381249",
							"author_id": "783214",
							"text": "Twitter API v2: Early Access releasednnToday we announced Early Access to the first endpoints of the new Twitter API!nn#TwitterAPI #EarlyAccess #VersionBump https://t.co/g7v3aeIbtQ"
						  }
						],
						"includes": {
						  "users": [
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
							}
						  ]
						},
						"meta": {
							"result_count": 3
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
				spaceID: "1DXxyRYNejbKM",
				opts: SpaceTweetsLookupOpts{
					Expansions: []Expansion{ExpansionAuthorID},
					UserFields: []UserField{UserFieldCreatedAt, UserFieldDescription},
				},
			},
			want: &SpaceTweetsLookupResponse{
				Raw: &TweetRaw{
					Tweets: []*TweetObj{
						{
							ID:       "1389270063807598594",
							AuthorID: "1065249714214457345",
							Text:     "now, everyone with 600 or more followers can host a Space.nnbased on what we've learned, these accounts are likely to have a good experience hosting because of their existing audience. before bringing the ability to create a Space to everyone, we're focused on a few things. :thread:",
						},
						{
							ID:       "1354143047324299264",
							AuthorID: "783214",
							Text:     "Academics are one of the biggest groups using the #TwitterAPI to research what's happening. Their work helps make the world (&amp; Twitter) a better place, and now more than ever, we must enable more of it. nIntroducing :drum_with_drumsticks: the Academic Research product track!nhttps://t.co/nOFiGewAV2",
						},
						{
							ID:       "1293595870563381249",
							AuthorID: "783214",
							Text:     "Twitter API v2: Early Access releasednnToday we announced Early Access to the first endpoints of the new Twitter API!nn#TwitterAPI #EarlyAccess #VersionBump https://t.co/g7v3aeIbtQ",
						},
					},
					Includes: &TweetRawIncludes{
						Users: []*UserObj{
							{
								ID:            "1065249714214457345",
								CreatedAt:     "2018-11-21T14:24:58.000Z",
								Name:          "Spaces",
								PinnedTweetID: "1389270063807598594",
								Description:   "Twitter Spaces is where live audio conversations happen.",
								UserName:      "TwitterSpaces",
							},
							{
								ID:          "783214",
								CreatedAt:   "2007-02-20T14:35:54.000Z",
								Name:        "Twitter",
								Description: "What's happening?!",
								UserName:    "Twitter",
							},
						},
					},
				},
				Meta: &SpaceTweetsLookupMeta{
					ResultCount: 3,
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
			got, err := c.SpaceTweetsLookup(context.Background(), tt.args.spaceID, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SpaceTweetsLookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.SpaceTweetsLookup() = %v, want %v", got, tt.want)
			}
		})
	}
}

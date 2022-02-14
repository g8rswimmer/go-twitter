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

func TestClient_ListLookup(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		listID string
		opts   ListLookupOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ListLookupResponse
		wantErr bool
	}{
		{
			name: "succes no options",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), listLookupEndpoint.urlID("", "list-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), listLookupEndpoint)
					}
					body := `{
						"data": {
						  "id": "84839422",
						  "name": "Official Twitter Accounts"
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
				listID: "list-1234",
			},
			want: &ListLookupResponse{
				Raw: &ListRaw{
					List: &ListObj{
						ID:   "84839422",
						Name: "Official Twitter Accounts",
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
			name: "succes options",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), listLookupEndpoint.urlID("", "list-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), listLookupEndpoint)
					}
					body := `{
						"data": {
						  "follower_count": 906,
						  "id": "84839422",
						  "name": "Official Twitter Accounts",
						  "owner_id": "783214"
						},
						"includes": {
						  "users": [
							{
							  "id": "783214",
							  "name": "Twitter",
							  "username": "Twitter"
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
				listID: "list-1234",
				opts: ListLookupOpts{
					Expansions: []Expansion{ExpansionOwnerID},
					ListFields: []ListField{ListFieldFollowerCount},
					UserFields: []UserField{UserFieldUserName},
				},
			},
			want: &ListLookupResponse{
				Raw: &ListRaw{
					List: &ListObj{
						ID:            "84839422",
						Name:          "Official Twitter Accounts",
						OwnerID:       "783214",
						FollowerCount: 906,
					},
					Includes: &ListRawIncludes{
						Users: []*UserObj{
							{
								ID:       "783214",
								Name:     "Twitter",
								UserName: "Twitter",
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
			got, err := c.ListLookup(context.Background(), tt.args.listID, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListLookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ListLookup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_UserListLookup(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		userID string
		opts   UserListLookupOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UserListLookupResponse
		wantErr bool
	}{
		{
			name: "success no options",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), userListLookupEndpoint.urlID("", "user-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userListLookupEndpoint)
					}
					body := `{
						"data": [
						  {
							"id": "1451305624956858369",
							"name": "Test List"
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
				userID: "user-1234",
			},
			want: &UserListLookupResponse{
				Raw: &UserListRaw{
					Lists: []*ListObj{
						{
							ID:   "1451305624956858369",
							Name: "Test List",
						},
					},
				},
				Meta: &UserListLookupMeta{
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
			name: "success with options",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), userListLookupEndpoint.urlID("", "user-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userListLookupEndpoint)
					}
					body := `{
						"data": [
						  {
							"follower_count": 0,
							"id": "1451305624956858369",
							"name": "Test List",
							"owner_id": "2244994945"
						  }
						],
						"includes": {
						  "users": [
							{
							  "username": "TwitterDev",
							  "id": "2244994945",
							  "created_at": "2013-12-14T04:35:55.000Z",
							  "name": "Twitter Dev"
							}
						  ]
						},
						"meta": {
						  "result_count": 1
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				userID: "user-1234",
				opts: UserListLookupOpts{
					Expansions: []Expansion{ExpansionOwnerID},
					ListFields: []ListField{ListFieldFollowerCount},
					UserFields: []UserField{UserFieldUserName},
				},
			},
			want: &UserListLookupResponse{
				Raw: &UserListRaw{
					Lists: []*ListObj{
						{
							ID:            "1451305624956858369",
							Name:          "Test List",
							FollowerCount: 0,
							OwnerID:       "2244994945",
						},
					},
					Includes: &ListRawIncludes{
						Users: []*UserObj{
							{
								ID:        "2244994945",
								Name:      "Twitter Dev",
								UserName:  "TwitterDev",
								CreatedAt: "2013-12-14T04:35:55.000Z",
							},
						},
					},
				},
				Meta: &UserListLookupMeta{
					ResultCount: 1,
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
			got, err := c.UserListLookup(context.Background(), tt.args.userID, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UserListLookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.UserListLookup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ListTweetLookup(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		listID string
		opts   ListTweetLookupOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ListTweetLookupResponse
		wantErr bool
	}{
		{
			name: "success no options",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), listTweetLookupEndpoint.urlID("", "list-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), listTweetLookupEndpoint)
					}
					body := `{
						"data": [
						  {
							"id": "1067094924124872705",
							"text": "Just getting started with Twitter APIs? Find out what you need in order to build an app. Watch this video! https://t.co/Hg8nkfoizN"
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
				listID: "list-1234",
			},
			want: &ListTweetLookupResponse{
				Raw: &TweetRaw{
					Tweets: []*TweetObj{
						{
							ID:   "1067094924124872705",
							Text: "Just getting started with Twitter APIs? Find out what you need in order to build an app. Watch this video! https://t.co/Hg8nkfoizN",
						},
					},
				},
				Meta: &ListTweetLookupMeta{
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
			name: "success with options",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), listTweetLookupEndpoint.urlID("", "list-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), listTweetLookupEndpoint)
					}
					body := `{
						"data": [
						  {
							"id": "1067094924124872705",
							"text": "Just getting started with Twitter APIs? Find out what you need in order to build an app. Watch this video! https://t.co/Hg8nkfoizN",
							"author_id": "2244994945"
						  }
						],
						"includes": {
							"users": [
							  {
								"verified": true,
								"username": "TwitterDev",
								"id": "2244994945",
								"name": "Twitter Dev"
							  }
							]
						  },
						"meta": {
						  "result_count": 1
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				listID: "list-1234",
				opts: ListTweetLookupOpts{
					Expansions: []Expansion{ExpansionAuthorID},
					UserFields: []UserField{UserFieldVerified},
				},
			},
			want: &ListTweetLookupResponse{
				Raw: &TweetRaw{
					Tweets: []*TweetObj{
						{
							ID:       "1067094924124872705",
							Text:     "Just getting started with Twitter APIs? Find out what you need in order to build an app. Watch this video! https://t.co/Hg8nkfoizN",
							AuthorID: "2244994945",
						},
					},
					Includes: &TweetRawIncludes{
						Users: []*UserObj{
							{
								Verified: true,
								UserName: "TwitterDev",
								ID:       "2244994945",
								Name:     "Twitter Dev",
							},
						},
					},
				},
				Meta: &ListTweetLookupMeta{
					ResultCount: 1,
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
			got, err := c.ListTweetLookup(context.Background(), tt.args.listID, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListTweetLookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ListTweetLookup() = %v, want %v", got, tt.want)
			}
		})
	}
}

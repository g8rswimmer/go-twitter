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

func TestClient_AddListMember(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		listID string
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ListAddMemberResponse
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodPost {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodPost)
					}
					if strings.Contains(req.URL.String(), listMemberEndpoint.urlID("", "list-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), listMemberEndpoint)
					}
					body := `{
						"data": {
						  "is_member": true
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
				userID: "user-1234",
			},
			want: &ListAddMemberResponse{
				List: &ListMemberData{
					Member: true,
				},
				RateLimit: &RateLimit{
					Limit:     15,
					Remaining: 12,
					Reset:     Epoch(1644461060),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := c.AddListMember(context.Background(), tt.args.listID, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.AddListMember() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.AddListMember() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_RemoveListMember(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		listID string
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ListRemoveMemberResponse
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodDelete {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodDelete)
					}
					if strings.Contains(req.URL.String(), listMemberEndpoint.urlID("", "list-1234")+"/user-1234") == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), listMemberEndpoint)
					}
					body := `{
						"data": {
						  "is_member": false
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
				userID: "user-1234",
			},
			want: &ListRemoveMemberResponse{
				List: &ListMemberData{
					Member: false,
				},
				RateLimit: &RateLimit{
					Limit:     15,
					Remaining: 12,
					Reset:     Epoch(1644461060),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := c.RemoveListMember(context.Background(), tt.args.listID, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.RemoveListMember() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.RemoveListMember() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ListUserMembers(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		listID string
		opts   ListUserMembersOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ListUserMembersResponse
		wantErr bool
	}{
		{
			name: "success  no options",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), listMemberEndpoint.urlID("", "list-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), listMemberEndpoint)
					}
					body := `{
						"data": [
						  {
							"id": "1065249714214457345",
							"name": "Spaces",
							"username": "TwitterSpaces"
						  }
						],
						"meta": {
						  "result_count": 1,
						  "next_token": "5676935732641845249"
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
			want: &ListUserMembersResponse{
				Raw: &UserRaw{
					Users: []*UserObj{
						{
							ID:       "1065249714214457345",
							Name:     "Spaces",
							UserName: "TwitterSpaces",
						},
					},
				},
				Meta: &ListUserMembersMeta{
					ResultCount: 1,
					NextToken:   "5676935732641845249",
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
			name: "success options",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), listMemberEndpoint.urlID("", "list-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), listMemberEndpoint)
					}
					body := `{
						"data": [
						  {
							"name": "Spaces",
							"id": "1065249714214457345",
							"username": "TwitterSpaces",
							"pinned_tweet_id": "1451239134798942208"
						  }
						],
						"includes": {
						  "tweets": [
							{
							  "id": "1451239134798942208",
							  "text": "the time has arrived -- we’re now rolling out the ability for everyone on iOS and Android to host a Spacennif this is your first time hosting, welcome! here’s a refresher on how https://t.co/cLH8z0bocy"
							}
						  ]
						},
						"meta": {
						  "result_count": 1,
						  "next_token": "5676935732641845249"
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
				opts: ListUserMembersOpts{
					Expansions: []Expansion{ExpansionPinnedTweetID},
					UserFields: []UserField{UserFieldUserName},
					MaxResults: 1,
				},
			},
			want: &ListUserMembersResponse{
				Raw: &UserRaw{
					Users: []*UserObj{
						{
							ID:            "1065249714214457345",
							Name:          "Spaces",
							UserName:      "TwitterSpaces",
							PinnedTweetID: "1451239134798942208",
						},
					},
					Includes: &UserRawIncludes{
						Tweets: []*TweetObj{
							{
								ID:   "1451239134798942208",
								Text: "the time has arrived -- we’re now rolling out the ability for everyone on iOS and Android to host a Spacennif this is your first time hosting, welcome! here’s a refresher on how https://t.co/cLH8z0bocy",
							},
						},
					},
				},
				Meta: &ListUserMembersMeta{
					ResultCount: 1,
					NextToken:   "5676935732641845249",
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
			got, err := c.ListUserMembers(context.Background(), tt.args.listID, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListUserMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ListUserMembers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_UserListMemberships(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		userID string
		opts   UserListMembershipsOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UserListMembershipsResponse
		wantErr bool
	}{
		{
			name: "success  no options",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), userListMemberEndpoint.urlID("", "user-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userListMemberEndpoint)
					}
					body := `{
						"data": [
						  {
							"id": "1450519480132509697",
							"name": "Twitter"
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
			want: &UserListMembershipsResponse{
				Raw: &UserListMembershipsRaw{
					Lists: []*ListObj{
						{
							ID:   "1450519480132509697",
							Name: "Twitter",
						},
					},
				},
				Meta: &UserListMembershipsMeta{
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
			name: "success  no options",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), userListMemberEndpoint.urlID("", "user-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userListMemberEndpoint)
					}
					body := `{
						"data": [
						  {
							"follower_count": 5,
							"id": "1451951974291689472",
							"name": "Twitter",
							"owner_id": "1227213680120479745"
						  }
						],
						"includes": {
						  "users": [
							{
							  "name": "구돆",
							  "created_at": "2020-02-11T12:52:11.000Z",
							  "id": "1227213680120479745",
							  "username": "Follow__Y0U"
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
			},
			want: &UserListMembershipsResponse{
				Raw: &UserListMembershipsRaw{
					Lists: []*ListObj{
						{
							ID:            "1451951974291689472",
							Name:          "Twitter",
							FollowerCount: 5,
							OwnerID:       "1227213680120479745",
						},
					},
					Includes: &ListRawIncludes{
						Users: []*UserObj{
							{
								ID:        "1227213680120479745",
								Name:      "구돆",
								UserName:  "Follow__Y0U",
								CreatedAt: "2020-02-11T12:52:11.000Z",
							},
						},
					},
				},
				Meta: &UserListMembershipsMeta{
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
			got, err := c.UserListMemberships(context.Background(), tt.args.userID, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UserListMemberships() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.UserListMemberships() = %v, want %v", got, tt.want)
			}
		})
	}
}

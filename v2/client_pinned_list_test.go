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

func TestClient_AddUserPinList(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		userID string
		listID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UserPinListResponse
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
					if strings.Contains(req.URL.String(), userPinnedListEndpoint.urlID("", "user-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userPinnedListEndpoint)
					}
					body := `{
						"data": {
						  "pinned": true
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
				listID: "list-1234",
			},
			want: &UserPinListResponse{
				List: &UserPinListData{
					Pinned: true,
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
			got, err := c.UserPinList(context.Background(), tt.args.userID, tt.args.listID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.AddUserPinList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.AddUserPinList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_RemoveUserPinList(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		userID string
		listID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UserUnpinListResponse
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
					if strings.Contains(req.URL.String(), userPinnedListEndpoint.urlID("", "user-1234")+"/list-1234") == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userPinnedListEndpoint)
					}
					body := `{
						"data": {
						  "pinned": false
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
				listID: "list-1234",
			},
			want: &UserUnpinListResponse{
				List: &UserPinListData{
					Pinned: false,
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
			got, err := c.UserUnpinList(context.Background(), tt.args.userID, tt.args.listID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.RemoveUserPinList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.RemoveUserPinList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_UserPinnedLists(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		userID string
		opts   UserPinnedListsOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UserPinnedListsResponse
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
					if strings.Contains(req.URL.String(), userPinnedListEndpoint.urlID("", "user-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userPinnedListEndpoint)
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
					}
				}),
			},
			args: args{
				userID: "user-1234",
			},
			want: &UserPinnedListsResponse{
				Raw: &UserPinnedListsRaw{
					Lists: []*ListObj{
						{
							ID:   "1451305624956858369",
							Name: "Test List",
						},
					},
				},
				Meta: &UserPinnedListsMeta{
					ResultCount: 1,
				},
			},
		},
		{
			name: "success",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), userPinnedListEndpoint.urlID("", "user-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userPinnedListEndpoint)
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
				opts: UserPinnedListsOpts{
					Expansions: []Expansion{ExpansionOwnerID},
					ListFields: []ListField{ListFieldFollowerCount},
					UserFields: []UserField{UserFieldCreatedAt},
				},
			},
			want: &UserPinnedListsResponse{
				Raw: &UserPinnedListsRaw{
					Lists: []*ListObj{
						{
							ID:            "1451305624956858369",
							Name:          "Test List",
							OwnerID:       "2244994945",
							FollowerCount: 0,
						},
					},
					Includes: &ListRawIncludes{
						Users: []*UserObj{
							{
								UserName:  "TwitterDev",
								ID:        "2244994945",
								CreatedAt: "2013-12-14T04:35:55.000Z",
								Name:      "Twitter Dev",
							},
						},
					},
				},
				Meta: &UserPinnedListsMeta{
					ResultCount: 1,
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
			got, err := c.UserPinnedLists(context.Background(), tt.args.userID, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UserPinnedLists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.UserPinnedLists() = %v, want %v", got, tt.want)
			}
		})
	}
}

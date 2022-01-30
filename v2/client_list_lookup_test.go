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

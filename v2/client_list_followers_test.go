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

func TestClient_UserFollowList(t *testing.T) {
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
		want    *UserFollowListResponse
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
					if strings.Contains(req.URL.String(), userFollowedListEndpoint.urlID("", "user-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userFollowedListEndpoint)
					}
					body := `{
						"data": {
						  "following": true
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
			want: &UserFollowListResponse{
				List: &UserFollowListData{
					Following: true,
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
			got, err := c.UserFollowList(context.Background(), tt.args.userID, tt.args.listID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UserFollowList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.UserFollowList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_UserUnfollowList(t *testing.T) {
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
		want    *UserUnfollowListResponse
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
					if strings.Contains(req.URL.String(), userFollowedListEndpoint.urlID("", "user-1234")+"/list-1234") == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userFollowedListEndpoint)
					}
					body := `{
						"data": {
						  "following": false
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
			want: &UserUnfollowListResponse{
				List: &UserFollowListData{
					Following: false,
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
			got, err := c.UserUnfollowList(context.Background(), tt.args.userID, tt.args.listID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UserUnfollowList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.UserUnfollowList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_UserFollowedLists(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		userID string
		opts   UserFollowedListsOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UserFollowedListsResponse
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
					if strings.Contains(req.URL.String(), userFollowedListEndpoint.urlID("", "user-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userFollowedListEndpoint)
					}
					body := `{
						"data": [
						  {
							"id": "1630685563471",
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
			want: &UserFollowedListsResponse{
				Raw: &UserFollowedListsRaw{
					Lists: []*ListObj{
						{
							ID:   "1630685563471",
							Name: "Test List",
						},
					},
				},
				Meta: &UserFollowedListsMeta{
					ResultCount: 1,
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
					if strings.Contains(req.URL.String(), userFollowedListEndpoint.urlID("", "user-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userFollowedListEndpoint)
					}
					body := `{
						"data": [
						  {
							"follower_count": 123,
							"id": "1630685563471",
							"name": "Test List",
							"owner_id": "1324848235714736129"
						  }
						],
						"includes": {
						  "users": [
							{
							  "username": "alanbenlee",
							  "id": "1324848235714736129",
							  "created_at": "2009-08-28T18:30:45.000Z",
							  "name": "Alan Lee"
							}
						  ]
						},
						"meta": {
						  "result_count": 1
						}
					  }
					  `
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				userID: "user-1234",
				opts: UserFollowedListsOpts{
					Expansions: []Expansion{ExpansionOwnerID},
					ListFields: []ListField{ListFieldFollowerCount},
					UserFields: []UserField{UserFieldUserName},
				},
			},
			want: &UserFollowedListsResponse{
				Raw: &UserFollowedListsRaw{
					Lists: []*ListObj{
						{
							ID:            "1630685563471",
							Name:          "Test List",
							OwnerID:       "1324848235714736129",
							FollowerCount: 123,
						},
					},
					Includes: &ListRawIncludes{
						Users: []*UserObj{
							{
								ID:        "1324848235714736129",
								UserName:  "alanbenlee",
								CreatedAt: "2009-08-28T18:30:45.000Z",
								Name:      "Alan Lee",
							},
						},
					},
				},
				Meta: &UserFollowedListsMeta{
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
			got, err := c.UserFollowedLists(context.Background(), tt.args.userID, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UserFollowedLists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.UserFollowedLists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ListUserFollowers(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		listID string
		opts   ListUserFollowersOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ListUserFollowersResponse
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
					if strings.Contains(req.URL.String(), listUserFollowersEndpoint.urlID("", "list-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), listUserFollowersEndpoint)
					}
					body := `{
						"data": [
						  {
							"id": "1420055293082415107",
							"name": "Bořek Šindelka(he/him)",
							"username": "JustBorek"
						  }
						],
						"meta": {
						  "result_count": 1,
						  "next_token": "1714209892546977900"
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
			want: &ListUserFollowersResponse{
				Raw: &UserRaw{
					Users: []*UserObj{
						{
							ID:       "1420055293082415107",
							Name:     "Bořek Šindelka(he/him)",
							UserName: "JustBorek",
						},
					},
				},
				Meta: &ListUserFollowersMeta{
					ResultCount: 1,
					NextToken:   "1714209892546977900",
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
					if strings.Contains(req.URL.String(), listUserFollowersEndpoint.urlID("", "list-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), listUserFollowersEndpoint)
					}
					body := `{
						"data": [
						  {
							"pinned_tweet_id": "1442182396523257861",  
							"id": "1420055293082415107",
							"name": "Bořek Šindelka(he/him)",
							"username": "JustBorek",
							"created_at": "2021-07-27T16:16:23.000Z"
						  }
						],
						"includes": {
							"tweets": [
							  {
								"created_at": "2021-09-26T17:40:52.000Z",
								"id": "1442182396523257861",
								"text": "Yes couple of days back nI want to kill my self I'm still here because of some amazing people please share this is important to talk about #mentalhealth @JustBorek #wheelchair #DisabilityTwitter #MedTwitter @heatherpsyd @Tweetinggoddess @NashaterS @msfatale @castleDD https://t.co/9hkSPV9NB1"
							  }
							]
						  },
						"meta": {
						  "result_count": 1,
						  "next_token": "1714209892546977900"
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
				opts: ListUserFollowersOpts{
					Expansions:  []Expansion{ExpansionPinnedTweetID},
					UserFields:  []UserField{UserFieldCreatedAt},
					TweetFields: []TweetField{TweetFieldCreatedAt},
				},
			},
			want: &ListUserFollowersResponse{
				Raw: &UserRaw{
					Users: []*UserObj{
						{
							ID:            "1420055293082415107",
							Name:          "Bořek Šindelka(he/him)",
							UserName:      "JustBorek",
							CreatedAt:     "2021-07-27T16:16:23.000Z",
							PinnedTweetID: "1442182396523257861",
						},
					},
					Includes: &UserRawIncludes{
						Tweets: []*TweetObj{
							{
								ID:        "1442182396523257861",
								Text:      "Yes couple of days back nI want to kill my self I'm still here because of some amazing people please share this is important to talk about #mentalhealth @JustBorek #wheelchair #DisabilityTwitter #MedTwitter @heatherpsyd @Tweetinggoddess @NashaterS @msfatale @castleDD https://t.co/9hkSPV9NB1",
								CreatedAt: "2021-09-26T17:40:52.000Z",
							},
						},
					},
				},
				Meta: &ListUserFollowersMeta{
					ResultCount: 1,
					NextToken:   "1714209892546977900",
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
			got, err := c.ListUserFollowers(context.Background(), tt.args.listID, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListUserFollowers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ListUserFollowers() = %v, want %v", got, tt.want)
			}
		})
	}
}

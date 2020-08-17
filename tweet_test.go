package twitter

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestTweetLookupParameters_encode(t *testing.T) {
	type fields struct {
		ids         []string
		Expansions  []Expansion
		MediaFields []MediaField
		PlaceFields []PlaceField
		PollFields  []PollField
		TweetFields []TweetField
		UserFields  []UserField
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   url.Values
	}{
		{
			name:   "no queries",
			fields: fields{},
			args: args{
				req: httptest.NewRequest(http.MethodGet, "https://www.go-twitter.com", nil),
			},
			want: url.Values{},
		},
		{
			name: "queries",
			fields: fields{
				ids:         []string{"123", "678"},
				Expansions:  []Expansion{ExpansionAuthorID},
				MediaFields: []MediaField{MediaFieldType, MediaFieldWidth},
				PlaceFields: []PlaceField{PlaceFieldID, PlaceFieldPlaceType},
				PollFields:  []PollField{PollFieldOptions},
				TweetFields: []TweetField{TweetFieldPossiblySensitve, TweetFieldNonPublicMetrics},
				UserFields:  []UserField{UserFieldProfileImageURL, UserFieldUserName},
			},
			args: args{
				req: httptest.NewRequest(http.MethodGet, "https://www.go-twitter.com", nil),
			},
			want: url.Values{
				"ids":          []string{"123,678"},
				"expansions":   []string{strings.Join(expansionStringArray([]Expansion{ExpansionAuthorID}), ",")},
				"media.fields": []string{strings.Join(mediaFieldStringArray([]MediaField{MediaFieldType, MediaFieldWidth}), ",")},
				"place.fields": []string{strings.Join(placeFieldStringArray([]PlaceField{PlaceFieldID, PlaceFieldPlaceType}), ",")},
				"poll.fields":  []string{strings.Join(pollFieldStringArray([]PollField{PollFieldOptions}), ",")},
				"tweet.fields": []string{strings.Join(tweetFieldStringArray([]TweetField{TweetFieldPossiblySensitve, TweetFieldNonPublicMetrics}), ",")},
				"user.fields":  []string{strings.Join(userFieldStringArray([]UserField{UserFieldProfileImageURL, UserFieldUserName}), ",")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tlp := TweetLookupParameters{
				ids:         tt.fields.ids,
				Expansions:  tt.fields.Expansions,
				MediaFields: tt.fields.MediaFields,
				PlaceFields: tt.fields.PlaceFields,
				PollFields:  tt.fields.PollFields,
				TweetFields: tt.fields.TweetFields,
				UserFields:  tt.fields.UserFields,
			}
			tlp.encode(tt.args.req)
			if reflect.DeepEqual(tt.args.req.URL.Query(), tt.want) == false {
				t.Errorf("TweetLookupParameters.encode() got %v want %v", tt.args.req.URL.Query(), tt.want)
			}
		})
	}
}

func TestTweet_Lookup(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		ids        []string
		parameters TweetLookupParameters
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           TweetLookups
		wantErr        bool
		wantTweetError *TweetErrorResponse
	}{
		{
			name: "success id",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
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
				parameters: TweetLookupParameters{
					UserFields: []UserField{UserFieldVerified, UserFieldUserName, UserFieldID, UserFieldName},
				},
			},
			want: TweetLookups{
				"1067094924124872705": TweetLookup{
					Tweet: TweetObj{
						AuthorID:  "2244994945",
						CreatedAt: "2018-11-26T16:37:10.000Z",
						Text:      "Just getting started with Twitter APIs? Find out what you need in order to build an app. Watch this video! https://t.co/Hg8nkfoizN",
						ID:        "1067094924124872705",
					},
					User: &UserObj{
						Verified: true,
						UserName: "TwitterDev",
						ID:       "2244994945",
						Name:     "Twitter Dev",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "success ids",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
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
				parameters: TweetLookupParameters{
					UserFields: []UserField{UserFieldVerified, UserFieldUserName, UserFieldID, UserFieldName},
				},
			},
			want: TweetLookups{
				"1261326399320715264": TweetLookup{
					Tweet: TweetObj{
						AuthorID:  "2244994945",
						CreatedAt: "2020-05-15T16:03:42.000Z",
						Text:      "Tune in to the @MongoDB @Twitch stream featuring our very own @suhemparack to learn about Twitter Developer Labs - starting now! https://t.co/fAWpYi3o5O",
						ID:        "1261326399320715264",
					},
					User: &UserObj{
						Verified: true,
						UserName: "TwitterDev",
						ID:       "2244994945",
						Name:     "Twitter Dev",
					},
				},
				"1278347468690915330": TweetLookup{
					Tweet: TweetObj{
						AuthorID:  "783214",
						CreatedAt: "2020-07-01T15:19:21.000Z",
						Text:      "Good news and bad news: \n\n2020 is half over",
						ID:        "1278347468690915330",
					},
					User: &UserObj{
						Verified: true,
						UserName: "Twitter",
						ID:       "783214",
						Name:     "Twitter",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "tweet error",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					body := `{
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
				parameters: TweetLookupParameters{
					UserFields: []UserField{UserFieldVerified, UserFieldUserName, UserFieldID, UserFieldName},
				},
			},
			want:    nil,
			wantErr: true,
			wantTweetError: &TweetErrorResponse{
				StatusCode: http.StatusBadRequest,
				Title:      "Invalid Request",
				Detail:     "One or more parameters to your request was invalid.",
				Type:       "https://api.twitter.com/2/problems/invalid-request",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tl := &Tweet{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := tl.Lookup(context.Background(), tt.args.ids, tt.args.parameters)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tweet.Lookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tweet.Lookup() = %v, want %v", got, tt.want)
			}
			var tweetErr *TweetErrorResponse
			if errors.As(err, &tweetErr) {
				if !reflect.DeepEqual(tweetErr, tt.wantTweetError) {
					t.Errorf("Tweet.Lookup() Error = %+v, want %+v", tweetErr, tt.wantTweetError)
				}
			}
		})
	}
}

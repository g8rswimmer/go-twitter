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

func TestClient_UserRetweet(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		userID  string
		tweetID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UserRetweetResponse
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodPost {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodPost)
					}
					if strings.Contains(req.URL.String(), userManageRetweetEndpoint.urlID("", "user-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userManageRetweetEndpoint)
					}
					body := `{
						"data": {
						  "retweeted": true
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				userID:  "user-1234",
				tweetID: "tweet-1234",
			},
			want: &UserRetweetResponse{
				Data: &RetweetData{
					Retweeted: true,
				},
			},
			wantErr: false,
		},
		{
			name:   "No User ID",
			fields: fields{},
			args: args{
				tweetID: "tweet-1234",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "No Tweet ID",
			fields: fields{},
			args: args{
				userID: "user-1234",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := c.UserRetweet(context.Background(), tt.args.userID, tt.args.tweetID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UserRetweet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.UserRetweet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_DeleteUserRetweet(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		userID  string
		tweetID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *DeleteUserRetweetResponse
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodDelete {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodDelete)
					}
					if strings.Contains(req.URL.String(), userManageRetweetEndpoint.urlID("", "user-1234")+"/tweet-1234") == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userManageRetweetEndpoint)
					}
					body := `{
							"data": {
							  "retweeted": false
							}
						  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				userID:  "user-1234",
				tweetID: "tweet-1234",
			},
			want: &DeleteUserRetweetResponse{
				Data: &RetweetData{
					Retweeted: false,
				},
			},
			wantErr: false,
		},
		{
			name:   "No User ID",
			fields: fields{},
			args: args{
				tweetID: "tweet-1234",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "No Tweet ID",
			fields: fields{},
			args: args{
				userID: "user-1234",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := c.DeleteUserRetweet(context.Background(), tt.args.userID, tt.args.tweetID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteUserRetweet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.DeleteUserRetweet() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func TestClient_CreateTweet(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		tweet CreateTweetRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CreateTweetResponse
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
					if strings.Contains(req.URL.String(), string(tweetCreateEndpoint)) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetCreateEndpoint)
					}
					body := `{
						"data": {
						  "id": "1445880548472328192",
						  "text": "Hello world!"
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusCreated,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				tweet: CreateTweetRequest{
					Text: "Hello world!",
				},
			},
			want: &CreateTweetResponse{
				Tweet: &CreateTweetData{
					Text: "Hello world!",
					ID:   "1445880548472328192",
				},
			},
			wantErr: false,
		},
		{
			name:   "Invalid Request",
			fields: fields{},
			args: args{
				tweet: CreateTweetRequest{},
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
			got, err := c.CreateTweet(context.Background(), tt.args.tweet)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateTweet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.CreateTweet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_DeleteTweet(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *DeleteTweetResponse
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
					if strings.Contains(req.URL.String(), tweetDeleteEndpoint.urlID("", "1445880548472328192")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetDeleteEndpoint)
					}
					body := `{
						"data": {
						  "deleted": true
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				id: "1445880548472328192",
			},
			want: &DeleteTweetResponse{
				Tweet: &DeleteTweetData{
					Deleted: true,
				},
			},
			wantErr: false,
		},
		{
			name:    "No ID",
			fields:  fields{},
			args:    args{},
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
			got, err := c.DeleteTweet(context.Background(), tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteTweet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.DeleteTweet() = %v, want %v", got, tt.want)
			}
		})
	}
}

package twitter

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

func TestClient_TweetRecentSearch(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		query string
		opts  TweetRecentSearchOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TweetRecentSearchResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := c.TweetRecentSearch(context.Background(), tt.args.query, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.TweetRecentSearch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.TweetRecentSearch() = %v, want %v", got, tt.want)
			}
		})
	}
}

package twitter

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestUserFieldOptions_addQuery(t *testing.T) {
	type fields struct {
		Expansions  []Expansion
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
				Expansions:  []Expansion{ExpansionAuthorID},
				TweetFields: []TweetField{TweetFieldPossiblySensitve, TweetFieldNonPublicMetrics},
				UserFields:  []UserField{UserFieldProfileImageURL, UserFieldUserName},
			},
			args: args{
				req: httptest.NewRequest(http.MethodGet, "https://www.go-twitter.com", nil),
			},
			want: url.Values{
				"expansions":   []string{strings.Join(expansionStringArray([]Expansion{ExpansionAuthorID}), ",")},
				"tweet.fields": []string{strings.Join(tweetFieldStringArray([]TweetField{TweetFieldPossiblySensitve, TweetFieldNonPublicMetrics}), ",")},
				"user.fields":  []string{strings.Join(userFieldStringArray([]UserField{UserFieldProfileImageURL, UserFieldUserName}), ",")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UserFieldOptions{
				Expansions:  tt.fields.Expansions,
				TweetFields: tt.fields.TweetFields,
				UserFields:  tt.fields.UserFields,
			}
			u.addQuery(tt.args.req)
			if reflect.DeepEqual(tt.args.req.URL.Query(), tt.want) == false {
				t.Errorf("UserFieldOptions.addQuery() got %v want %v", tt.args.req.URL.Query(), tt.want)
			}
		})
	}
}

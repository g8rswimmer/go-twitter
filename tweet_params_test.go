package twitter

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestTweetOptions_addQuery(t *testing.T) {
	type fields struct {
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
			opts := TweetFieldOptions{
				Expansions:  tt.fields.Expansions,
				MediaFields: tt.fields.MediaFields,
				PlaceFields: tt.fields.PlaceFields,
				PollFields:  tt.fields.PollFields,
				TweetFields: tt.fields.TweetFields,
				UserFields:  tt.fields.UserFields,
			}
			opts.addQuery(tt.args.req)
			if reflect.DeepEqual(tt.args.req.URL.Query(), tt.want) == false {
				t.Errorf("TweetOptions.addQuery() got %v want %v", tt.args.req.URL.Query(), tt.want)
			}
		})
	}
}

func TestTweetRecentSearchOptions_addQuery(t *testing.T) {
	type fields struct {
		query     string
		StartTime time.Time
		EndTime   time.Time
		MaxResult int
		NextToken string
		SinceID   string
		UntilID   string
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
			name: "no queries",
			fields: fields{
				query: "python",
			},
			args: args{
				req: httptest.NewRequest(http.MethodGet, "https://www.go-twitter.com", nil),
			},
			want: url.Values{
				"query": []string{"python"},
			},
		},
		{
			name: "queries",
			fields: fields{
				query:     "python",
				NextToken: "112233445566",
				StartTime: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2020, time.February, 20, 0, 0, 0, 0, time.UTC),
			},
			args: args{
				req: httptest.NewRequest(http.MethodGet, "https://www.go-twitter.com", nil),
			},
			want: url.Values{
				"query":      []string{"python"},
				"next_token": []string{"112233445566"},
				"end_time":   []string{"2020-02-20T00:00:00Z"},
				"start_time": []string{"2020-01-01T00:00:00Z"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := TweetRecentSearchOptions{
				query:     tt.fields.query,
				StartTime: tt.fields.StartTime,
				EndTime:   tt.fields.EndTime,
				MaxResult: tt.fields.MaxResult,
				NextToken: tt.fields.NextToken,
				SinceID:   tt.fields.SinceID,
				UntilID:   tt.fields.UntilID,
			}
			opts.addQuery(tt.args.req)
			if reflect.DeepEqual(tt.args.req.URL.Query(), tt.want) == false {
				t.Errorf("TweetRecentSearchOptions.addQuery() got %v want %v", tt.args.req.URL.Query(), tt.want)
			}
		})
	}
}

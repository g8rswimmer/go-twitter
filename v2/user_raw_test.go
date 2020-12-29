package twitter

import (
	"reflect"
	"testing"
)

func TestUserRaw_UserDictionaries(t *testing.T) {
	type fields struct {
		Users        []*UserObj
		Includes     *UserRawIncludes
		Errors       []*ErrorObj
		dictionaries map[string]*UserDictionary
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]*UserDictionary
	}{
		{
			name: "success",
			fields: fields{
				Users: []*UserObj{
					{
						ID:            "2244994945",
						Name:          "Twitter Dev",
						UserName:      "TwitterDev",
						CreatedAt:     "2013-12-14T04:35:55.000Z",
						PinnedTweetID: "1255542774432063488",
					},
					{
						ID:            "783214",
						Name:          "Twitter",
						UserName:      "Twitter",
						CreatedAt:     "2007-02-20T14:35:54.000Z",
						PinnedTweetID: "1274087687469715457",
					},
				},
				Includes: &UserRawIncludes{
					Tweets: []*TweetObj{
						{
							ID:        "1255542774432063488",
							CreatedAt: "2020-04-29T17:01:38.000Z",
							Text:      "During these unprecedented times, what’s happening on Twitter can help the world better understand &amp; respond to the pandemic. \n\nWe're launching a free COVID-19 stream endpoint so qualified devs &amp; researchers can study the public conversation in real-time. https://t.co/BPqMcQzhId",
						},
						{
							ID:        "1274087687469715457",
							CreatedAt: "2020-06-19T21:12:30.000Z",
							Text:      "📍 Minneapolis\n🗣️ @FredTJoseph https://t.co/lNTOkyguG1",
						},
					},
				},
			},
			want: map[string]*UserDictionary{
				"2244994945": {
					User: UserObj{
						ID:            "2244994945",
						Name:          "Twitter Dev",
						UserName:      "TwitterDev",
						CreatedAt:     "2013-12-14T04:35:55.000Z",
						PinnedTweetID: "1255542774432063488",
					},
					PinnedTweet: &TweetObj{
						ID:        "1255542774432063488",
						CreatedAt: "2020-04-29T17:01:38.000Z",
						Text:      "During these unprecedented times, what’s happening on Twitter can help the world better understand &amp; respond to the pandemic. \n\nWe're launching a free COVID-19 stream endpoint so qualified devs &amp; researchers can study the public conversation in real-time. https://t.co/BPqMcQzhId",
					},
				},
				"783214": {
					User: UserObj{
						ID:            "783214",
						Name:          "Twitter",
						UserName:      "Twitter",
						CreatedAt:     "2007-02-20T14:35:54.000Z",
						PinnedTweetID: "1274087687469715457",
					},
					PinnedTweet: &TweetObj{
						ID:        "1274087687469715457",
						CreatedAt: "2020-06-19T21:12:30.000Z",
						Text:      "📍 Minneapolis\n🗣️ @FredTJoseph https://t.co/lNTOkyguG1",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRaw{
				Users:        tt.fields.Users,
				Includes:     tt.fields.Includes,
				Errors:       tt.fields.Errors,
				dictionaries: tt.fields.dictionaries,
			}
			if got := u.UserDictionaries(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRaw.UserDictionaries() = %v, want %v", got, tt.want)
			}
		})
	}
}

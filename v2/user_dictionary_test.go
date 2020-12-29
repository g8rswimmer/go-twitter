package twitter

import (
	"reflect"
	"testing"
)

func TestCreateUserDictionary(t *testing.T) {
	type args struct {
		user     UserObj
		includes *UserRawIncludes
	}
	tests := []struct {
		name string
		args args
		want *UserDictionary
	}{
		{
			name: "success",
			args: args{
				user: UserObj{
					ID:            "2244994945",
					Name:          "Twitter Dev",
					UserName:      "TwitterDev",
					CreatedAt:     "2013-12-14T04:35:55.000Z",
					PinnedTweetID: "1255542774432063488",
				},
				includes: &UserRawIncludes{
					Tweets: []*TweetObj{
						{
							ID:        "1255542774432063488",
							CreatedAt: "2020-04-29T17:01:38.000Z",
							Text:      "During these unprecedented times, what’s happening on Twitter can help the world better understand &amp; respond to the pandemic. \n\nWe're launching a free COVID-19 stream endpoint so qualified devs &amp; researchers can study the public conversation in real-time. https://t.co/BPqMcQzhId",
						},
					},
				},
			},
			want: &UserDictionary{
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateUserDictionary(tt.args.user, tt.args.includes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateUserDictionary() = %v, want %v", got, tt.want)
			}
		})
	}
}

package twitter

import (
	"reflect"
	"testing"
)

func TestTweetRaw_TweetDictionaries(t *testing.T) {
	type fields struct {
		Tweets       []*TweetObj
		Includes     *TweetRawIncludes
		Errors       []*ErrorObj
		dictionaries map[string]*TweetDictionary
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]*TweetDictionary
	}{
		{
			name: "success",
			fields: fields{
				Tweets: []*TweetObj{
					{
						Text:      "Tune in to the @MongoDB @Twitch stream featuring our very own @suhemparack to learn about Twitter Developer Labs - starting now! https://t.co/fAWpYi3o5O",
						ID:        "1261326399320715264",
						AuthorID:  "2244994945",
						CreatedAt: "2020-05-15T16:03:42.000Z",
					},
					{
						Text:      "Good news and bad news: \n\n2020 is half over",
						ID:        "1278347468690915330",
						AuthorID:  "783214",
						CreatedAt: "2020-07-01T15:19:21.000Z",
					},
				},
				Includes: &TweetRawIncludes{
					Users: []*UserObj{
						{
							ID:       "2244994945",
							Verified: true,
							UserName: "TwitterDev",
							Name:     "Twitter Dev",
						},
						{
							ID:       "783214",
							Verified: true,
							UserName: "Twitter",
							Name:     "Twitter",
						},
					},
				},
			},
			want: map[string]*TweetDictionary{
				"1261326399320715264": {
					Tweet: TweetObj{
						Text:      "Tune in to the @MongoDB @Twitch stream featuring our very own @suhemparack to learn about Twitter Developer Labs - starting now! https://t.co/fAWpYi3o5O",
						ID:        "1261326399320715264",
						AuthorID:  "2244994945",
						CreatedAt: "2020-05-15T16:03:42.000Z",
					},
					Author: &UserObj{
						ID:       "2244994945",
						Verified: true,
						UserName: "TwitterDev",
						Name:     "Twitter Dev",
					},
					AttachmentMedia:  []*MediaObj{},
					AttachmentPolls:  []*PollObj{},
					Mentions:         []*TweetMention{},
					ReferencedTweets: []*TweetReference{},
				},
				"1278347468690915330": {
					Tweet: TweetObj{
						Text:      "Good news and bad news: \n\n2020 is half over",
						ID:        "1278347468690915330",
						AuthorID:  "783214",
						CreatedAt: "2020-07-01T15:19:21.000Z",
					},
					Author: &UserObj{
						ID:       "783214",
						Verified: true,
						UserName: "Twitter",
						Name:     "Twitter",
					},
					AttachmentMedia:  []*MediaObj{},
					AttachmentPolls:  []*PollObj{},
					Mentions:         []*TweetMention{},
					ReferencedTweets: []*TweetReference{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := &TweetRaw{
				Tweets:       tt.fields.Tweets,
				Includes:     tt.fields.Includes,
				Errors:       tt.fields.Errors,
				dictionaries: tt.fields.dictionaries,
			}
			if got := raw.TweetDictionaries(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TweetRaw.TweetDictionaries() = %v, want %v", got, tt.want)
			}
		})
	}
}

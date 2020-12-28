package twitter

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestCreateTweetDictionary(t *testing.T) {
	type args struct {
		tweet    TweetObj
		includes *TweetRawIncludes
	}
	tests := []struct {
		name string
		args args
		want *TweetDictionary
	}{
		{
			name: "complete dictionary",
			args: args{
				tweet: TweetObj{
					ID:   "1261326399320715264",
					Text: "Tune in to the @MongoDB @Twitch stream featuring our very own @suhemparack to learn about Twitter Developer Labs - starting now! https://t.co/fAWpYi3o5O",
					Attachments: &TweetAttachmentsObj{
						PollIDs:   []string{"1199786642468413448"},
						MediaKeys: []string{"13_1263145212760805376"},
					},
					Geo: &TweetGeoObj{
						PlaceID: "01a9a39529b27f36",
					},
					AuthorID:        "2244994945",
					InReplyToUserID: "783214",
					Entities: &EntitiesObj{
						Mentions: []EntityMentionObj{
							{
								EntityObj: EntityObj{
									Start: 15,
									End:   23,
								},
								UserName: "MongoDB",
							},
							{
								EntityObj: EntityObj{
									Start: 24,
									End:   31,
								},
								UserName: "Twitch",
							},
							{
								EntityObj: EntityObj{
									Start: 62,
									End:   74,
								},
								UserName: "suhemparack",
							},
						},
					},
					ReferencedTweets: []*TweetReferencedTweetObj{
						{
							Type: "quoted",
							ID:   "1261091720801980419",
						},
					},
				},
				includes: &TweetRawIncludes{
					Users: []*UserObj{
						{
							ID:       "2244994945",
							Name:     "Twitter Dev",
							UserName: "TwitterDev",
						},
						{
							Name:     "Twitter",
							ID:       "783214",
							UserName: "Twitter",
						},
						{
							Name:     "MongoDB",
							ID:       "18080585",
							UserName: "MongoDB",
						},
						{
							Name:     "Twitch",
							ID:       "309366491",
							UserName: "Twitch",
						},
						{
							Name:     "Suhem Parack",
							ID:       "857699969263964161",
							UserName: "suhemparack",
						},
					},
					Polls: []*PollObj{
						{
							ID:              "1199786642468413448",
							VotingStatus:    "closed",
							DurationMinutes: 1440,
							Options: []*PollOptionObj{
								{
									Position: 1,
									Label:    "C Sharp",
									Votes:    795,
								},
								{
									Position: 2,
									Label:    "C Hashtag",
									Votes:    156,
								},
							},
							EndDateTime: "2019-11-28T20:26:41.000Z",
						},
					},
					Media: []*MediaObj{
						{
							DurationMS: 46947,
							Type:       "video",
							Height:     1080,
							Key:        "13_1263145212760805376",
							PublicMetrics: &MediaMetricsObj{
								Views: 6909260,
							},
							PreviewImageURL: "https://pbs.twimg.com/media/EYeX7akWsAIP1_1.jpg",
							Width:           1920,
						},
					},
					Places: []*PlaceObj{
						{
							Geo: &PlaceGeoObj{
								Type: "Feature",
								BBox: []float64{
									-74.026675,
									40.683935,
									-73.910408,
									40.877483,
								},
								Properties: map[string]interface{}{},
							},
							CountryCode: "US",
							Name:        "Manhattan",
							ID:          "01a9a39529b27f36",
							PlaceType:   "city",
							Country:     "United States",
							FullName:    "Manhattan, NY",
						},
					},
					Tweets: []*TweetObj{
						{
							ID:       "1261091720801980419",
							AuthorID: "18080585",
							Text:     "Tomorrow (May 15) at 12pm EST (9am PST, 6pm CET), join us for a Twitch stream with @KukicAdo from MongoDB and @suhemparack from @TwitterDev! \n\nLearn about the new Twitter Developer Labs and how to get the most out of the new API with MongoDB: https://t.co/YbrbVNJrPe https://t.co/Oe4bMVpPmh",
						},
					},
				},
			},
			want: &TweetDictionary{
				Tweet: TweetObj{
					ID:   "1261326399320715264",
					Text: "Tune in to the @MongoDB @Twitch stream featuring our very own @suhemparack to learn about Twitter Developer Labs - starting now! https://t.co/fAWpYi3o5O",
					Attachments: &TweetAttachmentsObj{
						PollIDs:   []string{"1199786642468413448"},
						MediaKeys: []string{"13_1263145212760805376"},
					},
					Geo: &TweetGeoObj{
						PlaceID: "01a9a39529b27f36",
					},
					AuthorID:        "2244994945",
					InReplyToUserID: "783214",
					Entities: &EntitiesObj{
						Mentions: []EntityMentionObj{
							{
								EntityObj: EntityObj{
									Start: 15,
									End:   23,
								},
								UserName: "MongoDB",
							},
							{
								EntityObj: EntityObj{
									Start: 24,
									End:   31,
								},
								UserName: "Twitch",
							},
							{
								EntityObj: EntityObj{
									Start: 62,
									End:   74,
								},
								UserName: "suhemparack",
							},
						},
					},
					ReferencedTweets: []*TweetReferencedTweetObj{
						{
							Type: "quoted",
							ID:   "1261091720801980419",
						},
					},
				},
				Author: &UserObj{
					ID:       "2244994945",
					Name:     "Twitter Dev",
					UserName: "TwitterDev",
				},
				InReplyUser: &UserObj{
					Name:     "Twitter",
					ID:       "783214",
					UserName: "Twitter",
				},
				Mentions: []*TweetMention{
					{
						Mention: &EntityMentionObj{
							EntityObj: EntityObj{
								Start: 15,
								End:   23,
							},
							UserName: "MongoDB",
						},
						User: &UserObj{
							Name:     "MongoDB",
							ID:       "18080585",
							UserName: "MongoDB",
						},
					},
					{
						Mention: &EntityMentionObj{
							EntityObj: EntityObj{
								Start: 24,
								End:   31,
							},
							UserName: "Twitch",
						},
						User: &UserObj{
							Name:     "Twitch",
							ID:       "309366491",
							UserName: "Twitch",
						},
					},
					{
						Mention: &EntityMentionObj{
							EntityObj: EntityObj{
								Start: 62,
								End:   74,
							},
							UserName: "suhemparack",
						},
						User: &UserObj{
							Name:     "Suhem Parack",
							ID:       "857699969263964161",
							UserName: "suhemparack",
						},
					},
				},
				Place: &PlaceObj{
					Geo: &PlaceGeoObj{
						Type: "Feature",
						BBox: []float64{
							-74.026675,
							40.683935,
							-73.910408,
							40.877483,
						},
						Properties: map[string]interface{}{},
					},
					CountryCode: "US",
					Name:        "Manhattan",
					ID:          "01a9a39529b27f36",
					PlaceType:   "city",
					Country:     "United States",
					FullName:    "Manhattan, NY",
				},
				AttachmentPolls: []*PollObj{
					{
						ID:              "1199786642468413448",
						VotingStatus:    "closed",
						DurationMinutes: 1440,
						Options: []*PollOptionObj{
							{
								Position: 1,
								Label:    "C Sharp",
								Votes:    795,
							},
							{
								Position: 2,
								Label:    "C Hashtag",
								Votes:    156,
							},
						},
						EndDateTime: "2019-11-28T20:26:41.000Z",
					},
				},
				AttachmentMedia: []*MediaObj{
					{
						DurationMS: 46947,
						Type:       "video",
						Height:     1080,
						Key:        "13_1263145212760805376",
						PublicMetrics: &MediaMetricsObj{
							Views: 6909260,
						},
						PreviewImageURL: "https://pbs.twimg.com/media/EYeX7akWsAIP1_1.jpg",
						Width:           1920,
					},
				},
				ReferencedTweets: []*TweetReference{
					{
						Reference: &TweetReferencedTweetObj{
							Type: "quoted",
							ID:   "1261091720801980419",
						},
						TweetDictionary: &TweetDictionary{
							Tweet: TweetObj{
								ID:       "1261091720801980419",
								AuthorID: "18080585",
								Text:     "Tomorrow (May 15) at 12pm EST (9am PST, 6pm CET), join us for a Twitch stream with @KukicAdo from MongoDB and @suhemparack from @TwitterDev! \n\nLearn about the new Twitter Developer Labs and how to get the most out of the new API with MongoDB: https://t.co/YbrbVNJrPe https://t.co/Oe4bMVpPmh",
							},
							Author: &UserObj{
								Name:     "MongoDB",
								ID:       "18080585",
								UserName: "MongoDB",
							},
							AttachmentMedia:  []*MediaObj{},
							AttachmentPolls:  []*PollObj{},
							Mentions:         []*TweetMention{},
							ReferencedTweets: []*TweetReference{},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateTweetDictionary(tt.args.tweet, tt.args.includes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateTweetDictionary() = %+v, want %+v", got, tt.want)

				fmt.Println(func() string {
					enc, _ := json.MarshalIndent(got, "", "    ")
					return string(enc)
				}())
				fmt.Println(func() string {
					enc, _ := json.MarshalIndent(tt.want, "", "    ")
					return string(enc)
				}())
			}
		})
	}
}

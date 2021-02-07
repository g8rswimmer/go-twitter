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

func TestClient_UserMentionTimeline(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		userID string
		opts   UserMentionTimelineOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UserMentionTimelineResponse
		wantErr bool
	}{
		{
			name: "Success - Single ID Default",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), userMentionTimelineEndpoint.urlID("", "63046977")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), userTweetTimelineEndpoint)
					}
					body := `{
						"data": [
						  {
							"id": "1338980844036349953",
							"text": "@zeemacphee It’s an absolute honor to be the runner up to @happycamper great job all! 👏👏"
						  },
						  {
							"id": "1338973983312637955",
							"text": "I hope you enjoy your ENORMOUS grand prize @happycamper ‼️ https://t.co/KV48MENmBw https://t.co/oQg4MWW13a"
						  }
						],
						"meta": {
						  "oldest_id": "1336004278725513223",
						  "newest_id": "1338980844036349953",
						  "result_count": 2,
						  "next_token": "7140dibdnow9c7btw3w29kzu0unnfqs1lzcdi6s0vvj8z"
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				userID: "63046977",
			},
			want: &UserMentionTimelineResponse{
				Raw: &TweetRaw{
					Tweets: []*TweetObj{
						{
							ID:   "1338980844036349953",
							Text: "@zeemacphee It’s an absolute honor to be the runner up to @happycamper great job all! 👏👏",
						},
						{
							ID:   "1338973983312637955",
							Text: "I hope you enjoy your ENORMOUS grand prize @happycamper ‼️ https://t.co/KV48MENmBw https://t.co/oQg4MWW13a",
						},
					},
				},
				Meta: &UserTimelineMeta{
					ResultCount: 2,
					OldestID:    "1336004278725513223",
					NewestID:    "1338980844036349953",
					NextToken:   "7140dibdnow9c7btw3w29kzu0unnfqs1lzcdi6s0vvj8z",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := c.UserMentionTimeline(context.Background(), tt.args.userID, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UserMentionTimeline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.UserMentionTimeline() = %v, want %v", got, tt.want)
			}
		})
	}
}

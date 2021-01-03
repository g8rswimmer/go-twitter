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
		{
			name: "success",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), string(tweetRecentSearchEndpoint)) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetRecentSearchEndpoint)
					}
					body := `{
						"data": [
						  {
							"id": "1279990139888918528",
							"text": "Python now online for you !!\n\nWith the advent and acceptance of AI, Robotics, Python has become an inevitable factor in software development industry and most looked out skill both Nationally and Internationally. \n\nCoupon code: GVUP9\nCall: 9482303905/9482163905 https://t.co/ZFXCDJedAh"
						  },
						  {
							"id": "1279990133463429120",
							"text": "RT @McQubit: Building Neural Networks with Python Code and Math in Detail — II https://t.co/l6PKTTFGkv #machine_learning #programming #math…"
						  },
						  {
							"id": "1279990118355476480",
							"text": "RT @SunnyVaram: Top 10 Natural Language Processing Online Courses https://t.co/oAGqkHdS8H via @https://twitter.com/analyticsinme \n#DataScie…"
						  },
						  {
							"id": "1279990114584875009",
							"text": "RT @mohitnihalani7: LINK IN BIO......\n\n#programming #coding #programmer #developer #python #code #technology #coder #javascript #java #comp…"
						  },
						  {
							"id": "1279990108968665088",
							"text": "RT @mohitnihalani7: LINK IN BIO......\n\n#programming #coding #programmer #developer #python #code #technology #coder #javascript #java #comp…"
						  },
						  {
							"id": "1279990090828320769",
							"text": "RT @SunnyVaram: Top 10 Natural Language Processing Online Courses https://t.co/oAGqkHdS8H via @https://twitter.com/analyticsinme \n#DataScie…"
						  },
						  {
							"id": "1279990084398387201",
							"text": "RT @mohitnihalani7: LINK IN BIO......\n\n#programming #coding #programmer #developer #python #code #technology #coder #javascript #java #comp…"
						  },
						  {
							"id": "1279990076748038145",
							"text": "RT @gp_pulipaka: Best Machine Learning and Data Science #Books 2020. #BigData #Analytics #DataScience #IoT #IIoT #PyTorch #Python #RStats #…"
						  },
						  {
							"id": "1279990069105917952",
							"text": "RT @SunnyVaram: Top 10 Natural Language Processing Online Courses https://t.co/oAGqkHdS8H via @https://twitter.com/analyticsinme \n#DataScie…"
						  },
						  {
							"id": "1279990063888076800",
							"text": "RT @mohitnihalani7: LINK IN BIO......\n\n#programming #coding #programmer #developer #python #code #technology #coder #javascript #java #comp…"
						  }
						],
						"meta": {
						  "newest_id": "1279990139888918528",
						  "oldest_id": "1279990063888076800",
						  "result_count": 10,
						  "next_token": "b26v89c19zqg8o3fo7gghep0wmpt92c0wn0jiqwtc7tdp"
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				query: "phython",
			},
			want: &TweetRecentSearchResponse{
				Raw: &TweetRaw{
					Tweets: []*TweetObj{
						{
							ID:   "1279990139888918528",
							Text: "Python now online for you !!\n\nWith the advent and acceptance of AI, Robotics, Python has become an inevitable factor in software development industry and most looked out skill both Nationally and Internationally. \n\nCoupon code: GVUP9\nCall: 9482303905/9482163905 https://t.co/ZFXCDJedAh",
						},
						{
							ID:   "1279990133463429120",
							Text: "RT @McQubit: Building Neural Networks with Python Code and Math in Detail — II https://t.co/l6PKTTFGkv #machine_learning #programming #math…",
						},
						{
							ID:   "1279990118355476480",
							Text: "RT @SunnyVaram: Top 10 Natural Language Processing Online Courses https://t.co/oAGqkHdS8H via @https://twitter.com/analyticsinme \n#DataScie…",
						},
						{
							ID:   "1279990114584875009",
							Text: "RT @mohitnihalani7: LINK IN BIO......\n\n#programming #coding #programmer #developer #python #code #technology #coder #javascript #java #comp…",
						},
						{
							ID:   "1279990108968665088",
							Text: "RT @mohitnihalani7: LINK IN BIO......\n\n#programming #coding #programmer #developer #python #code #technology #coder #javascript #java #comp…",
						},
						{
							ID:   "1279990090828320769",
							Text: "RT @SunnyVaram: Top 10 Natural Language Processing Online Courses https://t.co/oAGqkHdS8H via @https://twitter.com/analyticsinme \n#DataScie…",
						},
						{
							ID:   "1279990084398387201",
							Text: "RT @mohitnihalani7: LINK IN BIO......\n\n#programming #coding #programmer #developer #python #code #technology #coder #javascript #java #comp…",
						},
						{
							ID:   "1279990076748038145",
							Text: "RT @gp_pulipaka: Best Machine Learning and Data Science #Books 2020. #BigData #Analytics #DataScience #IoT #IIoT #PyTorch #Python #RStats #…",
						},
						{
							ID:   "1279990069105917952",
							Text: "RT @SunnyVaram: Top 10 Natural Language Processing Online Courses https://t.co/oAGqkHdS8H via @https://twitter.com/analyticsinme \n#DataScie…",
						},
						{
							ID:   "1279990063888076800",
							Text: "RT @mohitnihalani7: LINK IN BIO......\n\n#programming #coding #programmer #developer #python #code #technology #coder #javascript #java #comp…",
						},
					},
				},
				Meta: &TweetRecentSearchMeta{
					NewestID:    "1279990139888918528",
					OldestID:    "1279990063888076800",
					ResultCount: 10,
					NextToken:   "b26v89c19zqg8o3fo7gghep0wmpt92c0wn0jiqwtc7tdp",
				},
			},
			wantErr: false,
		},
		{
			name: "success-optinoal",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), string(tweetRecentSearchEndpoint)) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetRecentSearchEndpoint)
					}
					body := `{
						"data": [
						  {
							"id": "1279990139888918528",
							"text": "Python now online for you !!\n\nWith the advent and acceptance of AI, Robotics, Python has become an inevitable factor in software development industry and most looked out skill both Nationally and Internationally. \n\nCoupon code: GVUP9\nCall: 9482303905/9482163905 https://t.co/ZFXCDJedAh",
							"lang": "en",
							"created_at": "2020-07-06T04:11:35.000Z",
							"conversation_id": "1279991360007929856"
						  },
						  {
							"id": "1279990133463429120",
							"text": "RT @McQubit: Building Neural Networks with Python Code and Math in Detail — II https://t.co/l6PKTTFGkv #machine_learning #programming #math…",
							"lang": "en",
							"created_at": "2020-07-06T04:11:34.000Z",
							"conversation_id": "1279991355326906369"
						  },
						  {
							"id": "1279990118355476480",
							"text": "RT @SunnyVaram: Top 10 Natural Language Processing Online Courses https://t.co/oAGqkHdS8H via @https://twitter.com/analyticsinme \n#DataScie…",
							"lang": "en",
							"created_at": "2020-07-06T04:11:34.000Z",
							"conversation_id": "1279991354223927296"
						  },
						  {
							"id": "1279990114584875009",
							"text": "RT @mohitnihalani7: LINK IN BIO......\n\n#programming #coding #programmer #developer #python #code #technology #coder #javascript #java #comp…",
							"lang": "en",
							"created_at": "2020-07-06T04:11:34.000Z",
							"conversation_id": "1279991354194624512"
						  },
						  {
							"id": "1279990108968665088",
							"text": "RT @mohitnihalani7: LINK IN BIO......\n\n#programming #coding #programmer #developer #python #code #technology #coder #javascript #java #comp…",
							"lang": "en",
							"created_at": "2020-07-06T04:11:29.000Z",
							"conversation_id": "1279991332421992448"
						  },
						  {
							"id": "1279990090828320769",
							"text": "RT @SunnyVaram: Top 10 Natural Language Processing Online Courses https://t.co/oAGqkHdS8H via @https://twitter.com/analyticsinme \n#DataScie…",
							"lang": "en",
							"created_at": "2020-07-06T04:11:27.000Z",
							"conversation_id": "1279991325237153794"
						  },
						  {
							"id": "1279990084398387201",
							"text": "RT @mohitnihalani7: LINK IN BIO......\n\n#programming #coding #programmer #developer #python #code #technology #coder #javascript #java #comp…",
							"lang": "en",
							"created_at": "2020-07-06T04:11:22.000Z",
							"conversation_id": "1279991304148172802"
						  },
						  {
							"id": "1279990076748038145",
							"text": "RT @gp_pulipaka: Best Machine Learning and Data Science #Books 2020. #BigData #Analytics #DataScience #IoT #IIoT #PyTorch #Python #RStats #…",
							"lang": "en",
							"created_at": "2020-07-06T04:11:21.000Z",
							"conversation_id": "1279991301249826816"
						  },
						  {
							"id": "1279990069105917952",
							"text": "RT @SunnyVaram: Top 10 Natural Language Processing Online Courses https://t.co/oAGqkHdS8H via @https://twitter.com/analyticsinme \n#DataScie…",
							"lang": "en",
							"created_at": "2020-07-06T04:11:21.000Z",
							"conversation_id": "1279991298443874305"
						  },
						  {
							"id": "1279990063888076800",
							"text": "RT @mohitnihalani7: LINK IN BIO......\n\n#programming #coding #programmer #developer #python #code #technology #coder #javascript #java #comp…",
							"lang": "en",
							"created_at": "2020-07-06T04:11:17.000Z",
							"conversation_id": "1279991282467815424"
						  }
						],
						"meta": {
						  "newest_id": "1279990139888918528",
						  "oldest_id": "1279990063888076800",
						  "result_count": 10,
						  "next_token": "b26v89c19zqg8o3fo7gghep0y5rnao6xpxi9raid7b0xp"
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				query: "phython",
				opts: TweetRecentSearchOpts{
					MaxResults:  10,
					TweetFields: []TweetField{TweetFieldCreatedAt, TweetFieldLanguage, TweetFieldConversationID},
				},
			},
			want: &TweetRecentSearchResponse{
				Raw: &TweetRaw{
					Tweets: []*TweetObj{
						{
							ID:             "1279990139888918528",
							Text:           "Python now online for you !!\n\nWith the advent and acceptance of AI, Robotics, Python has become an inevitable factor in software development industry and most looked out skill both Nationally and Internationally. \n\nCoupon code: GVUP9\nCall: 9482303905/9482163905 https://t.co/ZFXCDJedAh",
							Language:       "en",
							CreatedAt:      "2020-07-06T04:11:35.000Z",
							ConversationID: "1279991360007929856",
						},
						{
							ID:             "1279990133463429120",
							Text:           "RT @McQubit: Building Neural Networks with Python Code and Math in Detail — II https://t.co/l6PKTTFGkv #machine_learning #programming #math…",
							Language:       "en",
							CreatedAt:      "2020-07-06T04:11:34.000Z",
							ConversationID: "1279991355326906369",
						},
						{
							ID:             "1279990118355476480",
							Text:           "RT @SunnyVaram: Top 10 Natural Language Processing Online Courses https://t.co/oAGqkHdS8H via @https://twitter.com/analyticsinme \n#DataScie…",
							Language:       "en",
							CreatedAt:      "2020-07-06T04:11:34.000Z",
							ConversationID: "1279991354223927296",
						},
						{
							ID:             "1279990114584875009",
							Text:           "RT @mohitnihalani7: LINK IN BIO......\n\n#programming #coding #programmer #developer #python #code #technology #coder #javascript #java #comp…",
							Language:       "en",
							CreatedAt:      "2020-07-06T04:11:34.000Z",
							ConversationID: "1279991354194624512",
						},
						{
							ID:             "1279990108968665088",
							Text:           "RT @mohitnihalani7: LINK IN BIO......\n\n#programming #coding #programmer #developer #python #code #technology #coder #javascript #java #comp…",
							Language:       "en",
							CreatedAt:      "2020-07-06T04:11:29.000Z",
							ConversationID: "1279991332421992448",
						},
						{
							ID:             "1279990090828320769",
							Text:           "RT @SunnyVaram: Top 10 Natural Language Processing Online Courses https://t.co/oAGqkHdS8H via @https://twitter.com/analyticsinme \n#DataScie…",
							Language:       "en",
							CreatedAt:      "2020-07-06T04:11:27.000Z",
							ConversationID: "1279991325237153794",
						},
						{
							ID:             "1279990084398387201",
							Text:           "RT @mohitnihalani7: LINK IN BIO......\n\n#programming #coding #programmer #developer #python #code #technology #coder #javascript #java #comp…",
							Language:       "en",
							CreatedAt:      "2020-07-06T04:11:22.000Z",
							ConversationID: "1279991304148172802",
						},
						{
							ID:             "1279990076748038145",
							Text:           "RT @gp_pulipaka: Best Machine Learning and Data Science #Books 2020. #BigData #Analytics #DataScience #IoT #IIoT #PyTorch #Python #RStats #…",
							Language:       "en",
							CreatedAt:      "2020-07-06T04:11:21.000Z",
							ConversationID: "1279991301249826816",
						},
						{
							ID:             "1279990069105917952",
							Text:           "RT @SunnyVaram: Top 10 Natural Language Processing Online Courses https://t.co/oAGqkHdS8H via @https://twitter.com/analyticsinme \n#DataScie…",
							Language:       "en",
							CreatedAt:      "2020-07-06T04:11:21.000Z",
							ConversationID: "1279991298443874305",
						},
						{
							ID:             "1279990063888076800",
							Text:           "RT @mohitnihalani7: LINK IN BIO......\n\n#programming #coding #programmer #developer #python #code #technology #coder #javascript #java #comp…",
							Language:       "en",
							CreatedAt:      "2020-07-06T04:11:17.000Z",
							ConversationID: "1279991282467815424",
						},
					},
				},
				Meta: &TweetRecentSearchMeta{
					NewestID:    "1279990139888918528",
					OldestID:    "1279990063888076800",
					ResultCount: 10,
					NextToken:   "b26v89c19zqg8o3fo7gghep0y5rnao6xpxi9raid7b0xp",
				},
			},
			wantErr: false,
		},
		{
			name: "Bad Request",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), string(tweetLookupEndpoint)) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), tweetLookupEndpoint)
					}
					body := `{
						"errors": [
							{
								"parameters": {
									"id": [
										"aassd"
									]
								},
								"message": "The id query parameter value [aassd] does not match ^[0-9]{1,19}$"
							}
						],
						"title": "Invalid Request",
						"detail": "One or more parameters to your request was invalid.",
						"type": "https://api.twitter.com/2/problems/invalid-request"
					}`
					return &http.Response{
						StatusCode: http.StatusBadRequest,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				query: "nothing",
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

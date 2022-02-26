package twitter

import (
	"context"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestClient_SpacesSearch(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		query string
		opts  SpacesSearchOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *SpacesSearchResponse
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
					if strings.Contains(req.URL.String(), string(spaceSearchEndpoint)) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), spaceSearchEndpoint)
					}
					body := `{
						"data": [
						  {
							"host_ids": [
							  "2244994945"
							],
							"id": "1DXxyRYNejbKM",
							"state": "live",
							"title": "hello world 👋"
						  },
						  {
							"host_ids": [
							  "6253282"
							],
							"id": "1nAJELYEEPvGL",
							"state": "scheduled",
							"title": "Say hello to the Spaces endpoints"
						  }
						],
						"meta": {
						  "result_count": 2
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(body)),
						Header: func() http.Header {
							h := http.Header{}
							h.Add(rateLimit, "15")
							h.Add(rateRemaining, "12")
							h.Add(rateReset, "1644461060")
							return h
						}(),
					}
				}),
			},
			args: args{
				query: "hello",
				opts: SpacesSearchOpts{
					SpaceFields: []SpaceField{SpaceFieldTitle, SpaceFieldHostIDs},
				},
			},
			want: &SpacesSearchResponse{
				Raw: &SpacesRaw{
					Spaces: []*SpaceObj{
						{
							ID:    "1DXxyRYNejbKM",
							State: "live",
							Title: "hello world 👋",
							HostIDs: []string{
								"2244994945",
							},
						},
						{
							ID:    "1nAJELYEEPvGL",
							State: "scheduled",
							Title: "Say hello to the Spaces endpoints",
							HostIDs: []string{
								"6253282",
							},
						},
					},
				},
				Meta: &SpacesSearchMeta{
					ResultCount: 2,
				},
				RateLimit: &RateLimit{
					Limit:     15,
					Remaining: 12,
					Reset:     Epoch(1644461060),
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
			got, err := c.SpacesSearch(context.Background(), tt.args.query, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SpacesSearch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.SpacesSearch() = %v, want %v", got, tt.want)
			}
		})
	}
}

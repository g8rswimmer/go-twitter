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

func TestComplianceBatchJobObj_Upload(t *testing.T) {
	type fields struct {
		UploadURL string
		client    *http.Client
	}
	type args struct {
		ids io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				UploadURL: "https://wwww.test.com/update",
				client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodPut {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodPut)
					}
					if strings.Contains(req.URL.String(), "https://wwww.test.com/update") == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), "https://wwww.test.com/update")
					}
					body := ``
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
				ids: func() io.Reader {
					data := `1
					2
					3`
					return strings.NewReader(data)
				}(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ComplianceBatchJobObj{
				UploadURL: tt.fields.UploadURL,
				client:    tt.fields.client,
			}
			if err := c.Upload(context.Background(), tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("ComplianceBatchJobObj.Upload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestComplianceBatchJobObj_Download(t *testing.T) {
	type fields struct {
		DownloadURL string
		client      *http.Client
	}
	tests := []struct {
		name    string
		fields  fields
		want    *ComplianceBatchJobDownloadResponse
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				DownloadURL: "https://wwww.test.com/download",
				client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodGet {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodGet)
					}
					if strings.Contains(req.URL.String(), "https://wwww.test.com/download") == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), "https://wwww.test.com/download")
					}
					results := `{"id":"1265324480517361664","action":"delete","created_at":"2019-10-29T17:02:47.000Z","redacted_at":"2020-07-29T17:02:47.000Z","reason":"deleted"}`
					results += "\r\n"
					results += `{"id":"1263926741774581761","action":"delete","created_at":"2019-10-29T17:02:47.000Z","redacted_at":"2020-07-29T17:02:47.000Z","reason":"protected"}`
					results += "\r\n"
					results += `{"id":"1265324480517361669","action":"delete","created_at":"2019-10-29T17:02:47.000Z","redacted_at":"2020-07-29T17:02:47.000Z","reason":"suspended"}`

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(results)),
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
			want: &ComplianceBatchJobDownloadResponse{
				Results: []*ComplianceBatchJobResult{
					{
						ID:         "1265324480517361664",
						Action:     "delete",
						CreatedAt:  "2019-10-29T17:02:47.000Z",
						RedactedAt: "2020-07-29T17:02:47.000Z",
						Reason:     "deleted",
					},
					{
						ID:         "1263926741774581761",
						Action:     "delete",
						CreatedAt:  "2019-10-29T17:02:47.000Z",
						RedactedAt: "2020-07-29T17:02:47.000Z",
						Reason:     "protected",
					},
					{
						ID:         "1265324480517361669",
						Action:     "delete",
						CreatedAt:  "2019-10-29T17:02:47.000Z",
						RedactedAt: "2020-07-29T17:02:47.000Z",
						Reason:     "suspended",
					},
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
			c := ComplianceBatchJobObj{
				DownloadURL: tt.fields.DownloadURL,
				client:      tt.fields.client,
			}
			got, err := c.Download(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("ComplianceBatchJobObj.Download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ComplianceBatchJobObj.Download() = %v, want %v", got, tt.want)
			}
		})
	}
}

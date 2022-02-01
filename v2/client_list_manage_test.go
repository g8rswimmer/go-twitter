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

func TestClient_CreateList(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		list ListMetaData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ListCreateResponse
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodPost {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodPost)
					}
					if strings.Contains(req.URL.String(), listCreateEndpoint.url("")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), listCreateEndpoint)
					}
					body := `{
						"data": {
						  "id": "1441162269824405510",
						  "name": "test v2 create list"
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				list: ListMetaData{
					Name: func() *string {
						str := "test v2 create list"
						return &str
					}(),
				},
			},
			want: &ListCreateResponse{
				List: &ListCreateData{
					ID:   "1441162269824405510",
					Name: "test v2 create list",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := c.CreateList(context.Background(), tt.args.list)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.CreateList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_UpdateList(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		listID string
		update ListMetaData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ListUpdateResponse
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodPut {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodPut)
					}
					if strings.Contains(req.URL.String(), listUpdateEndpoint.urlID("", "list-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), listUpdateEndpoint)
					}
					body := `{
						"data": {
						  "updated": true
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				listID: "list-1234",
				update: ListMetaData{
					Name: func() *string {
						str := "test v2 create list"
						return &str
					}(),
				},
			},
			want: &ListUpdateResponse{
				List: &ListUpdateData{
					Updated: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := c.UpdateList(context.Background(), tt.args.listID, tt.args.update)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UpdateList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.UpdateList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_DeleteList(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		listID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ListDeleteResponse
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				Authorizer: &mockAuth{},
				Host:       "https://www.test.com",
				Client: mockHTTPClient(func(req *http.Request) *http.Response {
					if req.Method != http.MethodDelete {
						log.Panicf("the method is not correct %s %s", req.Method, http.MethodDelete)
					}
					if strings.Contains(req.URL.String(), listDeleteEndpoint.urlID("", "list-1234")) == false {
						log.Panicf("the url is not correct %s %s", req.URL.String(), listDeleteEndpoint)
					}
					body := `{
						"data": {
						  "deleted": true
						}
					  }`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(body)),
					}
				}),
			},
			args: args{
				listID: "list-1234",
			},
			want: &ListDeleteResponse{
				List: &ListDeleteData{
					Deleted: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := c.DeleteList(context.Background(), tt.args.listID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.DeleteList() = %v, want %v", got, tt.want)
			}
		})
	}
}

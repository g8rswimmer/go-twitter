package twitter

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func Test_rateFromHeader(t *testing.T) {
	type args struct {
		header http.Header
	}
	tests := []struct {
		name string
		args args
		want *RateLimit
	}{
		{
			name: "success",
			args: args{
				header: func() http.Header {
					h := http.Header{}
					h.Add(rateLimit, "15")
					h.Add(rateRemaining, "12")
					h.Add(rateReset, "1644461060")
					return h
				}(),
			},
			want: &RateLimit{
				Limit:     15,
				Remaining: 12,
				Reset:     Epoch(1644461060),
			},
		},
		{
			name: "fail",
			args: args{
				header: func() http.Header {
					h := http.Header{}
					h.Add(rateRemaining, "12")
					h.Add(rateReset, "1644461060")
					return h
				}(),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rateFromHeader(tt.args.header); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("rateFromHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRateLimitFromError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name  string
		args  args
		want  *RateLimit
		want1 bool
	}{
		{
			name: "error response",
			args: args{
				err: &ErrorResponse{
					RateLimit: &RateLimit{
						Limit:     15,
						Remaining: 12,
						Reset:     Epoch(1644461060),
					},
				},
			},
			want: &RateLimit{
				Limit:     15,
				Remaining: 12,
				Reset:     Epoch(1644461060),
			},
			want1: true,
		},
		{
			name: "http error",
			args: args{
				err: &HTTPError{
					RateLimit: &RateLimit{
						Limit:     15,
						Remaining: 12,
						Reset:     Epoch(1644461060),
					},
				},
			},
			want: &RateLimit{
				Limit:     15,
				Remaining: 12,
				Reset:     Epoch(1644461060),
			},
			want1: true,
		},
		{
			name: "response decode error",
			args: args{
				err: &ResponseDecodeError{
					RateLimit: &RateLimit{
						Limit:     15,
						Remaining: 12,
						Reset:     Epoch(1644461060),
					},
				},
			},
			want: &RateLimit{
				Limit:     15,
				Remaining: 12,
				Reset:     Epoch(1644461060),
			},
			want1: true,
		},
		{
			name: "error",
			args: args{
				err: errors.New("hit"),
			},
			want:  nil,
			want1: false,
		},
		{
			name:  "no error",
			args:  args{},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := RateLimitFromError(tt.args.err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RateLimitFromError() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("RateLimitFromError() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

package twitter

import (
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

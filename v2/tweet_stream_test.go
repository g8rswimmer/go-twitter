package twitter

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_StartTweetStreamMessage(t *testing.T) {
	type args struct {
		stream io.ReadCloser
	}
	tests := []struct {
		name string
		args args
		want []*TweetMessage
	}{
		{
			name: "tweet stream",
			args: args{
				stream: func() io.ReadCloser {
					stream := `{"data":{"id":"1","text":"hello"}}`
					stream += "\r\n"
					stream += `{"data":{"id":"2","text":"world"}}`
					stream += "\r\n"
					stream += `{"data":{"id":"3","text":"!!"}}`
					return io.NopCloser(strings.NewReader(stream))
				}(),
			},
			want: []*TweetMessage{
				{
					Raw: &TweetRaw{
						Tweets: []*TweetObj{
							{
								ID:   "1",
								Text: "hello",
							},
						},
					},
				},
				{
					Raw: &TweetRaw{
						Tweets: []*TweetObj{
							{
								ID:   "2",
								Text: "world",
							},
						},
					},
				},
				{
					Raw: &TweetRaw{
						Tweets: []*TweetObj{
							{
								ID:   "3",
								Text: "!!",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream := StartTweetStream(tt.args.stream)

			got := []*TweetMessage{}
			timer := time.NewTimer(time.Second * 5)

			func() {
				defer stream.Close()
				for {
					select {
					case msg := <-stream.Tweets():
						got = append(got, msg)
					case <-timer.C:
						return
					case err := <-stream.Err():
						t.Errorf("StartTweetStreamMessage error %v", err)
						return
					}
				}
			}()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StartTweetStreamMessage = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_StartTweetStreamSystem(t *testing.T) {
	type args struct {
		stream io.ReadCloser
	}
	tests := []struct {
		name string
		args args
		want []map[SystemMessageType]SystemMessage
	}{
		{
			name: "tweet stream",
			args: args{
				stream: func() io.ReadCloser {
					stream := `{"error":{"message":"Forced Disconnect: Too many connections. (Allowed Connections = 2)","sent":"2017-01-11T18:12:52+00:00"}}`
					stream += "\r\n"
					stream += `{"error":{"message":"Invalid date format for query parameter 'fromDate'. Expected format is 'yyyyMMddHHmm'. For example, '201701012315' for January 1st, 11:15 pm 2017 UTC.\n\n","sent":"2017-01-11T17:04:13+00:00"}}`
					stream += "\r\n"
					stream += `{"error":{"message":"Force closing connection to because it reached the maximum allowed backup (buffer size is ).","sent":"2017-01-11T17:04:13+00:00"}}`
					return io.NopCloser(strings.NewReader(stream))
				}(),
			},
			want: []map[SystemMessageType]SystemMessage{
				{
					ErrorMessageType: {
						Message: "Forced Disconnect: Too many connections. (Allowed Connections = 2)",
						Sent: func() time.Time {
							t, _ := time.Parse(time.RFC3339, "2017-01-11T18:12:52+00:00")
							return t
						}(),
					},
				},
				{
					ErrorMessageType: {
						Message: "Invalid date format for query parameter 'fromDate'. Expected format is 'yyyyMMddHHmm'. For example, '201701012315' for January 1st, 11:15 pm 2017 UTC.\n\n",
						Sent: func() time.Time {
							t, _ := time.Parse(time.RFC3339, "2017-01-11T17:04:13+00:00")
							return t
						}(),
					},
				},
				{
					ErrorMessageType: {
						Message: "Force closing connection to because it reached the maximum allowed backup (buffer size is ).",
						Sent: func() time.Time {
							t, _ := time.Parse(time.RFC3339, "2017-01-11T17:04:13+00:00")
							return t
						}(),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream := StartTweetStream(tt.args.stream)

			got := []map[SystemMessageType]SystemMessage{}
			timer := time.NewTimer(time.Second * 5)

			func() {
				defer stream.Close()
				for {
					select {
					case msg := <-stream.SystemMessages():
						got = append(got, msg)
					case <-timer.C:
						return
					case err := <-stream.Err():
						t.Errorf("StartTweetStreamMessage error %v", err)
						return
					}
				}
			}()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StartTweetStreamMessage = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_StartTweetStreamDisconnect(t *testing.T) {
	type args struct {
		stream io.ReadCloser
	}
	tests := []struct {
		name string
		args args
		want []*DisconnectionError
	}{
		{
			name: "disconnect stream",
			args: args{
				stream: func() io.ReadCloser {
					stream := `ewoJImVycm9ycyI6IFt7CgkJInRpdGxlIjogIm9wZXJhdGlvbmFsLWRpc2Nvbm5lY3QiLAoJCSJkaXNjb25uZWN0X3R5cGUiOiAiVXBzdHJlYW1PcGVyYXRpb25hbERpc2Nvbm5lY3QiLAoJCSJkZXRhaWwiOiAiVGhpcyBzdHJlYW0gaGFzIGJlZW4gZGlzY29ubmVjdGVkIHVwc3RyZWFtIGZvciBvcGVyYXRpb25hbCByZWFzb25zLiIsCgkJInR5cGUiOiAiaHR0cHM6Ly9hcGkudHdpdHRlci5jb20vMi9wcm9ibGVtcy9vcGVyYXRpb25hbC1kaXNjb25uZWN0IgoJfV0KfQ==`
					stream += "\r\n"
					stream += `ewoJInRpdGxlIjogIkNvbm5lY3Rpb25FeGNlcHRpb24iLAoJImRldGFpbCI6ICJUaGlzIHN0cmVhbSBpcyBjdXJyZW50bHkgYXQgdGhlIG1heGltdW0gYWxsb3dlZCBjb25uZWN0aW9uIGxpbWl0LiIsCgkiY29ubmVjdGlvbl9pc3N1ZSI6ICJUb29NYW55Q29ubmVjdGlvbnMiLAoJInR5cGUiOiAiaHR0cHM6Ly9hcGkudHdpdHRlci5jb20vMi9wcm9ibGVtcy9zdHJlYW1pbmctY29ubmVjdGlvbiIKfQ==`
					return io.NopCloser(strings.NewReader(stream))
				}(),
			},
			want: []*DisconnectionError{
				{
					Disconnections: []*Disconnection{
						{
							Title:          "operational-disconnect",
							DisconnectType: "UpstreamOperationalDisconnect",
							Detail:         "This stream has been disconnected upstream for operational reasons.",
							Type:           "https://api.twitter.com/2/problems/operational-disconnect",
						},
					},
					Connections: []*Connection{},
				},
				{
					Disconnections: []*Disconnection{},
					Connections: []*Connection{
						{
							Title:           "ConnectionException",
							ConnectionIssue: "TooManyConnections",
							Detail:          "This stream is currently at the maximum allowed connection limit.",
							Type:            "https://api.twitter.com/2/problems/streaming-connection",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream := StartTweetStream(tt.args.stream)

			got := []*DisconnectionError{}
			timer := time.NewTimer(time.Second * 5)

			func() {
				defer stream.Close()
				for {
					select {
					case msg := <-stream.DisconnectionError():
						got = append(got, msg)
					case <-timer.C:
						return
					case err := <-stream.Err():
						t.Errorf("Test_StartTweetStreamDisconnect error %v", err)
						return
					}
				}
			}()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Test_StartTweetStreamDisconnect = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_StartTweetStream(t *testing.T) {
	type args struct {
		stream io.ReadCloser
	}
	tests := []struct {
		name           string
		args           args
		wantSystem     []map[SystemMessageType]SystemMessage
		wantTweet      []*TweetMessage
		wantDisconnect []*DisconnectionError
	}{
		{
			name: "tweet stream",
			args: args{
				stream: func() io.ReadCloser {
					stream := `{"data":{"id":"1","text":"hello"}}`
					stream += "\r\n"
					stream += `{"error":{"message":"Forced Disconnect: Too many connections. (Allowed Connections = 2)","sent":"2017-01-11T18:12:52+00:00"}}`
					stream += "\r\n"
					stream += `{"data":{"id":"2","text":"world"}}`
					stream += "\r\n"
					stream += "\r\n"
					stream += "\r\n"
					stream += "\r\n"
					stream += `{"data":{"id":"3","text":"!!"}}`
					stream += "\r\n"
					stream += `{"error":{"message":"Invalid date format for query parameter 'fromDate'. Expected format is 'yyyyMMddHHmm'. For example, '201701012315' for January 1st, 11:15 pm 2017 UTC.\n\n","sent":"2017-01-11T17:04:13+00:00"}}`
					stream += "\r\n"
					stream += `{"error":{"message":"Force closing connection to because it reached the maximum allowed backup (buffer size is ).","sent":"2017-01-11T17:04:13+00:00"}}`
					stream += "\r\n"
					stream += `ewoJInRpdGxlIjogIkNvbm5lY3Rpb25FeGNlcHRpb24iLAoJImRldGFpbCI6ICJUaGlzIHN0cmVhbSBpcyBjdXJyZW50bHkgYXQgdGhlIG1heGltdW0gYWxsb3dlZCBjb25uZWN0aW9uIGxpbWl0LiIsCgkiY29ubmVjdGlvbl9pc3N1ZSI6ICJUb29NYW55Q29ubmVjdGlvbnMiLAoJInR5cGUiOiAiaHR0cHM6Ly9hcGkudHdpdHRlci5jb20vMi9wcm9ibGVtcy9zdHJlYW1pbmctY29ubmVjdGlvbiIKfQ==`
					return io.NopCloser(strings.NewReader(stream))
				}(),
			},
			wantTweet: []*TweetMessage{
				{
					Raw: &TweetRaw{
						Tweets: []*TweetObj{
							{
								ID:   "1",
								Text: "hello",
							},
						},
					},
				},
				{
					Raw: &TweetRaw{
						Tweets: []*TweetObj{
							{
								ID:   "2",
								Text: "world",
							},
						},
					},
				},
				{
					Raw: &TweetRaw{
						Tweets: []*TweetObj{
							{
								ID:   "3",
								Text: "!!",
							},
						},
					},
				},
			},
			wantSystem: []map[SystemMessageType]SystemMessage{
				{
					ErrorMessageType: {
						Message: "Forced Disconnect: Too many connections. (Allowed Connections = 2)",
						Sent: func() time.Time {
							t, _ := time.Parse(time.RFC3339, "2017-01-11T18:12:52+00:00")
							return t
						}(),
					},
				},
				{
					ErrorMessageType: {
						Message: "Invalid date format for query parameter 'fromDate'. Expected format is 'yyyyMMddHHmm'. For example, '201701012315' for January 1st, 11:15 pm 2017 UTC.\n\n",
						Sent: func() time.Time {
							t, _ := time.Parse(time.RFC3339, "2017-01-11T17:04:13+00:00")
							return t
						}(),
					},
				},
				{
					ErrorMessageType: {
						Message: "Force closing connection to because it reached the maximum allowed backup (buffer size is ).",
						Sent: func() time.Time {
							t, _ := time.Parse(time.RFC3339, "2017-01-11T17:04:13+00:00")
							return t
						}(),
					},
				},
			},
			wantDisconnect: []*DisconnectionError{
				{
					Disconnections: []*Disconnection{},
					Connections: []*Connection{
						{
							Title:           "ConnectionException",
							ConnectionIssue: "TooManyConnections",
							Detail:          "This stream is currently at the maximum allowed connection limit.",
							Type:            "https://api.twitter.com/2/problems/streaming-connection",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream := StartTweetStream(tt.args.stream)

			gotSystem := []map[SystemMessageType]SystemMessage{}
			gotTweet := []*TweetMessage{}
			gotDisconnect := []*DisconnectionError{}

			timer := time.NewTimer(time.Second * 5)

			func() {
				defer stream.Close()
				for {
					select {
					case sysMsg := <-stream.SystemMessages():
						gotSystem = append(gotSystem, sysMsg)
					case tweetMsg := <-stream.Tweets():
						gotTweet = append(gotTweet, tweetMsg)
					case disconnectMsg := <-stream.DisconnectionError():
						gotDisconnect = append(gotDisconnect, disconnectMsg)
					case <-timer.C:
						return
					case err := <-stream.Err():
						t.Errorf("StartTweetStreamMessage error %v", err)
						return
					}
				}
			}()

			if !reflect.DeepEqual(gotSystem, tt.wantSystem) {
				t.Errorf("StartTweetStreamMessage system= %v, want %v", gotSystem, tt.wantSystem)
			}

			if !reflect.DeepEqual(gotTweet, tt.wantTweet) {
				t.Errorf("StartTweetStreamMessage system= %v, want %v", gotTweet, tt.wantTweet)
			}

			if !reflect.DeepEqual(gotDisconnect, tt.wantDisconnect) {
				t.Errorf("StartTweetStreamMessage system= %v, want %v", gotDisconnect, tt.wantDisconnect)
			}

		})
	}
}

func Test_streamSeperator(t *testing.T) {
	type args struct {
		data  []byte
		atEOF bool
	}
	tests := []struct {
		name    string
		args    args
		want    int
		want1   []byte
		wantErr bool
	}{
		{
			name: "sperated",
			args: args{
				data: func() []byte {
					msg := `{"data":{"id":"1","text":"hello"}}`
					msg += "\r\n"
					msg += `{"data":{"id":"2","text":"world"}}`
					msg += "\r\n"
					msg += `{"data":{"id":"3","text":"!!"}}`
					return []byte(msg)
				}(),
				atEOF: false,
			},
			want:    len(`{"data":{"id":"1","text":"hello"}}`) + 2,
			want1:   []byte(`{"data":{"id":"1","text":"hello"}}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := streamSeparator(tt.args.data, tt.args.atEOF)
			if (err != nil) != tt.wantErr {
				t.Errorf("streamSeperator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("streamSeperator() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("streamSeperator() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestStreamError_Error(t *testing.T) {
	type fields struct {
		Type StreamErrorType
		Msg  string
		Err  error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "message with error",
			fields: fields{
				Type: TweetErrorType,
				Msg:  "test message",
				Err:  errors.New("wow"),
			},
			want: fmt.Sprintf("%s: test message wow", TweetErrorType),
		},
		{
			name: "message",
			fields: fields{
				Type: TweetErrorType,
				Msg:  "test message",
			},
			want: fmt.Sprintf("%s: test message", TweetErrorType),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := StreamError{
				Type: tt.fields.Type,
				Msg:  tt.fields.Msg,
				Err:  tt.fields.Err,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("StreamError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStreamError_Is(t *testing.T) {
	type fields struct {
		Type StreamErrorType
		Msg  string
		Err  error
	}
	type args struct {
		target error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "success",
			fields: fields{
				Type: TweetErrorType,
			},
			args: args{
				target: &StreamError{
					Type: TweetErrorType,
				},
			},
			want: true,
		},
		{
			name: "success wrapped",
			fields: fields{
				Type: TweetErrorType,
			},
			args: args{
				target: func() error {
					e := &StreamError{
						Type: TweetErrorType,
					}
					return fmt.Errorf("some error %w", e)
				}(),
			},
			want: true,
		},
		{
			name: "fail",
			fields: fields{
				Type: TweetErrorType,
			},
			args: args{
				target: &StreamError{
					Type: SystemErrorType,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &StreamError{
				Type: tt.fields.Type,
				Msg:  tt.fields.Msg,
				Err:  tt.fields.Err,
			}
			if got := errors.Is(tt.args.target, e); got != tt.want {
				t.Errorf("StreamError.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStreamError_Wrap(t *testing.T) {

	werr := errors.New("wow")

	e := &StreamError{
		Err: werr,
	}
	if !errors.Is(e, werr) {
		t.Error("want error unwrapped")
	}

}

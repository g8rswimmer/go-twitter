package twitter

import "context"



func TestClient_UserTweetTimeline(t *testing.T) {
	type fields struct {
		Authorizer Authorizer
		Client     *http.Client
		Host       string
	}
	type args struct {
		userID string
		opts   UserTweetTimelineOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UserTweetTimelineResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Authorizer: tt.fields.Authorizer,
				Client:     tt.fields.Client,
				Host:       tt.fields.Host,
			}
			got, err := c.UserTweetTimeline(context.Background(), tt.args.userID, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UserTweetTimeline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.UserTweetTimeline() = %v, want %v", got, tt.want)
			}
		})
	}
}

package twitter

import (
	"testing"
)

func TestCreateTweetReply_validate(t *testing.T) {
	type fields struct {
		ExcludeReplyUserIDs []string
		InReplyToTweetID    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				ExcludeReplyUserIDs: []string{"6253282"},
				InReplyToTweetID:    "1455953449422516226",
			},
			wantErr: false,
		},
		{
			name: "valid 2",
			fields: fields{
				InReplyToTweetID: "1455953449422516226",
			},
			wantErr: false,
		},
		{
			name: "invalid",
			fields: fields{
				ExcludeReplyUserIDs: []string{"6253282"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := CreateTweetReply{
				ExcludeReplyUserIDs: tt.fields.ExcludeReplyUserIDs,
				InReplyToTweetID:    tt.fields.InReplyToTweetID,
			}
			if err := r.validate(); (err != nil) != tt.wantErr {
				t.Errorf("CreateTweetReply.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateTweetPoll_validate(t *testing.T) {
	type fields struct {
		DurationMinutes int
		Options         []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				Options:         []string{"yes", "maybe", "no"},
				DurationMinutes: 120,
			},
			wantErr: false,
		},
		{
			name: "invalid",
			fields: fields{
				Options: []string{"yes", "maybe", "no"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := CreateTweetPoll{
				DurationMinutes: tt.fields.DurationMinutes,
				Options:         tt.fields.Options,
			}
			if err := p.validate(); (err != nil) != tt.wantErr {
				t.Errorf("CreateTweetPoll.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateTweetMedia_validate(t *testing.T) {
	type fields struct {
		IDs           []string
		TaggedUserIDs []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				IDs:           []string{"1455952740635586573"},
				TaggedUserIDs: []string{"2244994945", "6253282"},
			},
			wantErr: false,
		},
		{
			name: "valid 2",
			fields: fields{
				IDs: []string{"1455952740635586573"},
			},
			wantErr: false,
		},
		{
			name: "invalid",
			fields: fields{
				TaggedUserIDs: []string{"2244994945", "6253282"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := CreateTweetMedia{
				IDs:           tt.fields.IDs,
				TaggedUserIDs: tt.fields.TaggedUserIDs,
			}
			if err := m.validate(); (err != nil) != tt.wantErr {
				t.Errorf("CreateTweetMedia.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateTweetOps_validate(t *testing.T) {
	type fields struct {
		DirectMessageDeepLink string
		ForSuperFollowersOnly bool
		QuoteTweetID          string
		Text                  string
		ReplySettings         string
		Geo                   CreateTweetGeo
		Media                 CreateTweetMedia
		Poll                  CreateTweetPoll
		Reply                 CreateTweetReply
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				Text: "Hello World",
				Media: CreateTweetMedia{
					IDs: []string{"12345"},
				},
			},
			wantErr: false,
		},
		{
			name: "valid2",
			fields: fields{
				Media: CreateTweetMedia{
					IDs: []string{"12345"},
				},
			},
			wantErr: false,
		},
		{
			name: "valid3",
			fields: fields{
				Text: "Hello World",
			},
			wantErr: false,
		},
		{
			name:    "invalid",
			fields:  fields{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := CreateTweetRequest{
				DirectMessageDeepLink: tt.fields.DirectMessageDeepLink,
				ForSuperFollowersOnly: tt.fields.ForSuperFollowersOnly,
				QuoteTweetID:          tt.fields.QuoteTweetID,
				Text:                  tt.fields.Text,
				ReplySettings:         tt.fields.ReplySettings,
				Geo:                   tt.fields.Geo,
				Media:                 tt.fields.Media,
				Poll:                  tt.fields.Poll,
				Reply:                 tt.fields.Reply,
			}
			if err := opts.validate(); (err != nil) != tt.wantErr {
				t.Errorf("CreateTweetOps.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

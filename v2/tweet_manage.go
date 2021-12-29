package twitter

import "fmt"

type CreateTweetRequest struct {
	DirectMessageDeepLink string           `json:"direct_message_deep_link,omitempty"`
	ForSuperFollowersOnly bool             `json:"for_super_followers_only,omitempty"`
	QuoteTweetID          string           `json:"quote_tweet_id,omitempty"`
	Text                  string           `json:"text,omitemtpy"`
	ReplySettings         string           `json:"reply_settings"`
	Geo                   CreateTweetGeo   `json:"geo"`
	Media                 CreateTweetMedia `json:"media"`
	Poll                  CreateTweetPoll  `json:"poll"`
	Reply                 CreateTweetReply `json:"reply"`
}

func (t CreateTweetRequest) validate() error {
	if err := t.Media.validate(); err != nil {
		return fmt.Errorf("create tweet error: %w", err)
	}
	if err := t.Poll.validate(); err != nil {
		return fmt.Errorf("create tweet error: %w", err)
	}
	if err := t.Reply.validate(); err != nil {
		return fmt.Errorf("create tweet error: %w", err)
	}
	if len(t.Media.IDs) == 0 && len(t.Text) == 0 {
		return fmt.Errorf("create tweet text is required if no media ids %w", ErrParameter)
	}
	return nil
}

type CreateTweetGeo struct {
	PlaceID string `json:"place_id"`
}

type CreateTweetMedia struct {
	IDs           []string `json:"media_ids"`
	TaggedUserIDs []string `json:"tagged_user_ids"`
}

func (m CreateTweetMedia) validate() error {
	if len(m.TaggedUserIDs) > 0 && len(m.IDs) == 0 {
		return fmt.Errorf("media ids are required if taged user ids are present %w", ErrParameter)
	}
	return nil
}

type CreateTweetPoll struct {
	DurationMinutes int      `json:"duration_minutes"`
	Options         []string `json:"options"`
}

func (p CreateTweetPoll) validate() error {
	if len(p.Options) > 0 && p.DurationMinutes <= 0 {
		return fmt.Errorf("poll duration minutes are required with options %w", ErrParameter)
	}
	return nil
}

type CreateTweetReply struct {
	ExcludeReplyUserIDs []string `json:"exclude_reply_user_ids"`
	InReplyToTweetID    string   `json:"in_reply_to_tweet_id"`
}

func (r CreateTweetReply) validate() error {
	if len(r.ExcludeReplyUserIDs) > 0 && len(r.InReplyToTweetID) == 0 {
		return fmt.Errorf("reply in reply to tweet is needs to be present it exclude reply user ids are present %w", ErrParameter)
	}
	return nil
}

type CreateTweetData struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type CreateTweetResponse struct {
	Tweet *CreateTweetData `json:"data"`
}

type DeleteTweetData struct {
	Deleted bool `json:"deleted"`
}

type DeleteTweetResponse struct {
	Tweet *DeleteTweetData `json:"data"`
}

package model

// Media refers to any image, GIF, or video attached to a Tweet
type Media struct {
	Key              string       `json:"media_key"`
	Type             string       `json:"type"`
	DurationMS       int          `json:"duration_ms"`
	Height           int          `json:"height"`
	NonPublicMetrics MediaMetrics `json:"non_public_metrics"`
	OrganicMetrics   MediaMetrics `json:"organic_metrics"`
	PreviewImageURL  string       `json:"preview_image_url"`
	PromotedMetrics  MediaMetrics `json:"promoted_metrics"`
	PublicMetrics    MediaMetrics `json:"public_metrics"`
	Width            int          `json:"width"`
}

// MediaMetrics engagement metrics for the media content at the time of the request
type MediaMetrics struct {
	Playback0   int `json:"playback_0_count"`
	Playback100 int `json:"playback_100_count"`
	Playback25  int `json:"playback_25_count"`
	Playback50  int `json:"playback_50_count"`
	Playback75  int `json:"playback_75_count"`
	Views       int `json:"view_count"`
}

package model

// MediaField can expand the fields that are returned in the media object
type MediaField string

const (
	// MediaFieldDurationMS available when type is video. Duration in milliseconds of the video.
	MediaFieldDurationMS MediaField = "duration_ms"
	// MediaFieldHeight of this content in pixels.
	MediaFieldHeight MediaField = "height"
	// MediaFieldMediaKey unique identifier of the expanded media content.
	MediaFieldMediaKey MediaField = "media_key"
	// MediaFieldPreviewImageURL is the URL to the static placeholder preview of this content.
	MediaFieldPreviewImageURL MediaField = "preview_image_url"
	// MediaFieldType is the type of content (animated_gif, photo, video)
	MediaFieldType MediaField = "type"
	// MediaFieldURL is the URL of the content
	MediaFieldURL MediaField = "url"
	// MediaFieldWidth is the width of this content in pixels
	MediaFieldWidth MediaField = "width"
	// MediaFieldPublicMetrics is the public engagement metrics for the media content at the time of the request.
	MediaFieldPublicMetrics MediaField = "public_metrics"
	// MediaFieldNonPublicMetrics is the non-public engagement metrics for the media content at the time of the request.
	MediaFieldNonPublicMetrics MediaField = "non_public_metrics"
	// MediaFieldOrganicMetrics is the engagement metrics for the media content, tracked in an organic context, at the time of the request.
	MediaFieldOrganicMetrics MediaField = "organic_metrics"
	// MediaFieldPromotedMetrics is the URL to the static placeholder preview of this content.
	MediaFieldPromotedMetrics MediaField = "promoted_metrics"
)

// Media refers to any image, GIF, or video attached to a Tweet
type Media struct {
	Key              string       `json:"media_key"`
	Type             string       `json:"type"`
	URL              string       `json:"url"`
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

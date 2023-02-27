package twitter

// MediaField can expand the fields that are returned in the media object
type MediaField string

type MediaType string

// MediaType which is used to store values that can be any of two options, photo or video.
// This type can then be used to distinguish between the two media types in various contexts.
const (
	Photo MediaType = "photo"
	Video MediaType = "video"
)

// MediaPublicMetrics is the public engagement metrics for the media content at the time of the request.
// This includes engagement metrics tracked at the account level and engagement metrics tracked at the URL level.
// View count is the sum of view counts from both contexts.
type MediaPublicMetrics struct {
	ViewCount int64 `json:"view_count"`
}

// Variant is a video variant object that contains information about a specific video format.
// The variant with the highest bitrate is the format that is used when a video is played in the Twitter player.
// The other variants are provided for users who have slower connections or who have chosen to use data-saving mode.
type Variant struct {
	BitRate     *int64 `json:"bit_rate,omitempty"`
	ContentType string `json:"content_type"`
	URL         string `json:"url"`
}

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
	// MediaFieldVariants is the variants of the media content.
	MediaFieldVariants MediaField = "variants"
)

func mediaFieldStringArray(arr []MediaField) []string {
	strs := make([]string, len(arr))
	for i, field := range arr {
		strs[i] = string(field)
	}
	return strs
}

// MediaObj refers to any image, GIF, or video attached to a Tweet
type MediaObj struct {
	Key              string          `json:"media_key"`
	Type             string          `json:"type"`
	URL              string          `json:"url"`
	DurationMS       int             `json:"duration_ms"`
	Height           int             `json:"height"`
	NonPublicMetrics MediaMetricsObj `json:"non_public_metrics"`
	OrganicMetrics   MediaMetricsObj `json:"organic_metrics"`
	PreviewImageURL  string          `json:"preview_image_url"`
	PromotedMetrics  MediaMetricsObj `json:"promoted_metrics"`
	PublicMetrics    MediaMetricsObj `json:"public_metrics"`
	Width            int             `json:"width"`
	Variants         []Variant       `json:"variants,omitempty"`
}

// MediaMetricsObj engagement metrics for the media content at the time of the request
type MediaMetricsObj struct {
	Playback0   int `json:"playback_0_count"`
	Playback100 int `json:"playback_100_count"`
	Playback25  int `json:"playback_25_count"`
	Playback50  int `json:"playback_50_count"`
	Playback75  int `json:"playback_75_count"`
	Views       int `json:"view_count"`
}

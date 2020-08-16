package twitter

// Entities contains details about text that has a special meaning.
type Entities struct {
	Annotations []EntityAnnotation `json:"annotations"`
	URLs        []EntityURL        `json:"urls"`
	HashTags    []EntityTag        `json:"hashtags"`
	Mentions    []EntityMention    `json:"mentions"`
	CashTags    []EntityTag        `json:"cashtags"`
}

// Entity contains the start and end positions of the text
type Entity struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// EntityAnnotation contains details about annotations relative to the text.
type EntityAnnotation struct {
	Entity
	Probability    float64 `json:"probability"`
	Type           string  `json:"type"`
	NormalizedText string  `json:"normalized_text"`
}

// EntityURL contains details about text recognized as a URL.
type EntityURL struct {
	Entity
	URL         string `json:"url"`
	ExpandedURL string `json:"expanded_url"`
	DisplayURL  string `json:"display_url"`
	Status      string `json:"status"`
	Title       string `json:"title"`
	Desription  string `json:"description"`
	UnwoundURL  string `json:"unwound_url"`
}

// EntityTag contains details about text recognized as a tag
type EntityTag struct {
	Entity
	Tag string `json:"tag"`
}

// EntityMention contains details about text recognized as a user mention.
type EntityMention struct {
	Entity
	UserName string `json:"username"`
}

// WithHeld contains withholding details
type WithHeld struct {
	Copyright    bool     `json:"copyright"`
	CountryCodes []string `json:"country_codes"`
}

package twitter

// EntitiesObj contains details about text that has a special meaning.
type EntitiesObj struct {
	Annotations []EntityAnnotationObj `json:"annotations"`
	URLs        []EntityURLObj        `json:"urls"`
	HashTags    []EntityTagObj        `json:"hashtags"`
	Mentions    []EntityMentionObj    `json:"mentions"`
	CashTags    []EntityTagObj        `json:"cashtags"`
}

// EntityObj contains the start and end positions of the text
type EntityObj struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// EntityAnnotationObj contains details about annotations relative to the text.
type EntityAnnotationObj struct {
	EntityObj
	Probability    float64 `json:"probability"`
	Type           string  `json:"type"`
	NormalizedText string  `json:"normalized_text"`
}

// EntityURLObj contains details about text recognized as a URL.
type EntityURLObj struct {
	EntityObj
	URL         string `json:"url"`
	ExpandedURL string `json:"expanded_url"`
	DisplayURL  string `json:"display_url"`
	Status      int    `json:"status"`
	Title       string `json:"title"`
	Desription  string `json:"description"`
	UnwoundURL  string `json:"unwound_url"`
}

// EntityTagObj contains details about text recognized as a tag
type EntityTagObj struct {
	EntityObj
	Tag string `json:"tag"`
}

// EntityMentionObj contains details about text recognized as a user mention.
type EntityMentionObj struct {
	EntityObj
	UserName string `json:"username"`
}

// WithHeldObj contains withholding details
type WithHeldObj struct {
	Copyright    bool     `json:"copyright"`
	CountryCodes []string `json:"country_codes"`
}

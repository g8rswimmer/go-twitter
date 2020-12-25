package twitter

// PlaceField can expand the tweet primary object
type PlaceField string

const (
	// PlaceFieldContainedWithin returns the identifiers of known places that contain the referenced place.
	PlaceFieldContainedWithin PlaceField = "contained_within"
	// PlaceFieldCountry is the full-length name of the country this place belongs to.
	PlaceFieldCountry PlaceField = "country"
	// PlaceFieldCountryCode is the ISO Alpha-2 country code this place belongs to.
	PlaceFieldCountryCode PlaceField = "country_code"
	// PlaceFieldFullName is a longer-form detailed place name.
	PlaceFieldFullName PlaceField = "full_name"
	// PlaceFieldGeo contains place details in GeoJSON format.
	PlaceFieldGeo PlaceField = "geo"
	// PlaceFieldID is the unique identifier of the expanded place, if this is a point of interest tagged in the Tweet.
	PlaceFieldID PlaceField = "id"
	// PlaceFieldName is the short name of this place
	PlaceFieldName PlaceField = "name"
	// PlaceFieldPlaceType is specified the particular type of information represented by this place information, such as a city name, or a point of interest.
	PlaceFieldPlaceType PlaceField = "place_type"
)

func placeFieldStringArray(arr []PlaceField) []string {
	strs := make([]string, len(arr))
	for i, field := range arr {
		strs[i] = string(field)
	}
	return strs
}

// PlaceObj tagged in a Tweet is not a primary object on any endpoint
type PlaceObj struct {
	FullName        string       `json:"full_name,omitempty"`
	ID              string       `json:"id"`
	ContainedWithin []string     `json:"contained_within,omitempty"`
	Country         string       `json:"country,omitempty"`
	CountryCode     string       `json:"country_code,omitempty"`
	Geo             *PlaceGeoObj `json:"geo,omitempty"`
	Name            string       `json:"name"`
	PlaceType       string       `json:"place_type,omitempty"`
}

// PlaceGeoObj contains place details
type PlaceGeoObj struct {
	Type       string                 `json:"type"`
	BBox       []float64              `json:"bbox"`
	Properties map[string]interface{} `json:"properties"`
}

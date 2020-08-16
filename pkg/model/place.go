package model

// Place tagged in a Tweet is not a primary object on any endpoint
type Place struct {
	FullName        string   `json:"full_name"`
	ID              string   `json:"id"`
	ContainedWithin []string `json:"contained_within"`
	Country         string   `json:"country"`
	CountryCode     string   `json:"country_code"`
	Geo             PlaceGeo `json:"geo"`
	Name            string   `json:"name"`
	PlaceType       string   `json:"place_type"`
}

// PlaceGeo contains place details
type PlaceGeo struct {
	Type       string                 `json:"type"`
	BBox       []float64              `json:"bbox"`
	Properties map[string]interface{} `json:"properties"`
}

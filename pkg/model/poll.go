package model

// Poll included in a Tweet is not a primary object on any endpoint
type Poll struct {
	ID              string       `json:"id"`
	Options         []PollOption `json:"options"`
	DurationMinutes int          `json:"duration_minutes"`
	EndDateTime     string       `json:"end_datetime"`
	VotingStatus    string       `json:"voting_status"`
}

// PollOption contains objects describing each choice in the referenced poll.
type PollOption struct {
	Position int    `json:"position"`
	Label    string `json:"label"`
	Votes    int    `json:"votes"`
}

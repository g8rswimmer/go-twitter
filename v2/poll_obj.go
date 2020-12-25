package twitter

// PollField defines the fields of the expanded tweet
type PollField string

const (
	// PollFieldDurationMinutes specifies the total duration of this poll.
	PollFieldDurationMinutes PollField = "duration_minutes"
	// PollFieldEndDateTime specifies the end date and time for this poll.
	PollFieldEndDateTime PollField = "end_datetime"
	// PollFieldID is unique identifier of the expanded poll.
	PollFieldID PollField = "id"
	// PollFieldOptions contains objects describing each choice in the referenced poll.
	PollFieldOptions PollField = "options"
	// PollFieldVotingStatus indicates if this poll is still active and can receive votes, or if the voting is now closed.
	PollFieldVotingStatus PollField = "voting_status"
)

func pollFieldStringArray(arr []PollField) []string {
	strs := make([]string, len(arr))
	for i, field := range arr {
		strs[i] = string(field)
	}
	return strs
}

// PollObj included in a Tweet is not a primary object on any endpoint
type PollObj struct {
	ID              string           `json:"id"`
	Options         []*PollOptionObj `json:"options,omitempty"`
	DurationMinutes int              `json:"duration_minutes,omitempty"`
	EndDateTime     string           `json:"end_datetime,omitempty"`
	VotingStatus    string           `json:"voting_status,omitempty"`
}

// PollOptionObj contains objects describing each choice in the referenced poll.
type PollOptionObj struct {
	Position int    `json:"position"`
	Label    string `json:"label"`
	Votes    int    `json:"votes"`
}

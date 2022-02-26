package twitter

// SpaceField are the space field options
type SpaceField string

const (
	// SpaceFieldHostIDs is the space host ids field
	SpaceFieldHostIDs SpaceField = "host_ids"
	// SpaceFieldCreatedAt is the space created at field
	SpaceFieldCreatedAt SpaceField = "created_at"
	// SpaceFieldCreatorID is the space creator id field
	SpaceFieldCreatorID SpaceField = "creator_id"
	// SpaceFieldID is the space field id field
	SpaceFieldID SpaceField = "id"
	// SpaceFieldLang is the space language field
	SpaceFieldLang SpaceField = "lang"
	// SpaceFieldInvitedUserIDs is the space invited user ids field
	SpaceFieldInvitedUserIDs SpaceField = "invited_user_ids"
	// SpaceFieldParticipantCount is the space participant count field
	SpaceFieldParticipantCount SpaceField = "participant_count"
	// SpaceFieldSpeakerIDs is the space speaker ids field
	SpaceFieldSpeakerIDs SpaceField = "speaker_ids"
	// SpaceFieldStartedAt is the space started at field
	SpaceFieldStartedAt SpaceField = "started_at"
	// SpaceFieldEndedAt is the space ended at field
	SpaceFieldEndedAt SpaceField = "ended_at"
	// SpaceFieldSubscriberCount is the space subscriber count field
	SpaceFieldSubscriberCount SpaceField = "subscriber_count"
	// SpaceFieldTopicIDs is the space topic ids field
	SpaceFieldTopicIDs SpaceField = "topic_ids"
	// SpaceFieldState is the space state field
	SpaceFieldState SpaceField = "state"
	// SpaceFieldTitle is the space title field
	SpaceFieldTitle SpaceField = "title"
	// SpaceFieldUpdatedAt is the space updated at field
	SpaceFieldUpdatedAt SpaceField = "updated_at"
	// SpaceFieldScheduledStart is the space scheduled start field
	SpaceFieldScheduledStart SpaceField = "scheduled_start"
	// SpaceFieldTicketed is the space is ticketed field
	SpaceFieldTicketed SpaceField = "is_ticketed"
)

func spaceFieldStringArray(arr []SpaceField) []string {
	strs := make([]string, len(arr))
	for i, field := range arr {
		strs[i] = string(field)
	}
	return strs
}

// SpaceState is the enumeration of the space states
type SpaceState string

const (
	// SpaceStateAll is for all of the possible states
	SpaceStateAll SpaceState = "all"
	// SpaceStateLive is for only live states
	SpaceStateLive SpaceState = "live"
	// SpaceStateScheduled is for only scheduled states
	SpaceStateScheduled SpaceState = "scheduled"
)

// SpaceObj is the spaces object
type SpaceObj struct {
	ID               string   `json:"id"`
	State            string   `json:"state"`
	CreatedAt        string   `json:"created_at"`
	EndedAt          string   `json:"ended_at"`
	HostIDs          []string `json:"host_ids"`
	Lang             string   `json:"lang"`
	Ticketed         bool     `json:"is_ticketed"`
	InvitedUserIDs   []string `json:"invited_user_ids"`
	ParticipantCount int      `json:"participant_count"`
	ScheduledStart   string   `json:"scheduled_start"`
	SpeakerIDs       []string `json:"speaker_ids"`
	StartedAt        string   `json:"started_at"`
	Title            string   `json:"title"`
	TopicIDs         []string `json:"topic_ids"`
	UpdatedAt        string   `json:"updated_at"`
	CreatorID        string   `json:"creator_id"`
	SubscriberCount  int      `json:"subscriber_count"`
}

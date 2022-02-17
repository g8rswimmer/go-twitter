package twitter

type SpaceField string

const (
	SpaceFieldHostIDs          SpaceField = "host_ids"
	SpaceFieldCreatedAt        SpaceField = "created_at"
	SpaceFieldCreatorID        SpaceField = "creator_id"
	SpaceFieldID               SpaceField = "id"
	SpaceFieldLang             SpaceField = "lang"
	SpaceFieldInvitedUserIDs   SpaceField = "invited_user_ids"
	SpaceFieldParticipantCount SpaceField = "participant_count"
	SpaceFieldSpeakerIDs       SpaceField = "speaker_ids"
	SpaceFieldStartedAt        SpaceField = "started_at"
	SpaceFieldEndedAt          SpaceField = "ended_at"
	SpaceFieldSubscriberCount  SpaceField = "subscriber_count"
	SpaceFieldTopicIDs         SpaceField = "topic_ids"
	SpaceFieldState            SpaceField = "state"
	SpaceFieldTitle            SpaceField = "title"
	SpaceFieldUpdatedAt        SpaceField = "updated_at"
	SpaceFieldScheduledStart   SpaceField = "scheduled_start"
	SpaceFieldTicketed         SpaceField = "is_ticketed"
)

func spaceFieldStringArray(arr []SpaceField) []string {
	strs := make([]string, len(arr))
	for i, field := range arr {
		strs[i] = string(field)
	}
	return strs
}

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

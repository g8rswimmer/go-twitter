package twitter

// TopicField are the topic field options
type TopicField string

const (
	// TopicFieldID is the topic id field
	TopicFieldID TopicField = "id"
	// TopicFieldName is the topic name field
	TopicFieldName TopicField = "name"
	// TopicFieldDescription is the topic description field
	TopicFieldDescription TopicField = "description"
)

func topicFieldStringArray(arr []TopicField) []string {
	strs := make([]string, len(arr))
	for i, field := range arr {
		strs[i] = string(field)
	}
	return strs
}

// TopicObj is the topic object
type TopicObj struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

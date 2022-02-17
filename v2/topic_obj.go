package twitter

type TopicField string

const (
	TopicFieldID          TopicField = "id"
	TopicFieldName        TopicField = "name"
	TopicFieldDescription TopicField = "description"
)

func topocFieldStringArray(arr []TopicField) []string {
	strs := make([]string, len(arr))
	for i, field := range arr {
		strs[i] = string(field)
	}
	return strs
}

type TopicObj struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

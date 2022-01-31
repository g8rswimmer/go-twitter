package twitter

// ListField are the optional fields that can be included in the response
type ListField string

const (
	// ListFieldCreatedAt is the UTC datetime that the List was created on Twitter.
	ListFieldCreatedAt ListField = "created_at"
	// ListFieldFollowerCount shows how many users follow this List
	ListFieldFollowerCount ListField = "follower_count"
	// ListFieldMemberCount shows how many members are part of this List.
	ListFieldMemberCount ListField = "member_count"
	// ListFieldPrivate indicates if the list is private
	ListFieldPrivate ListField = "private"
	// ListFieldDescription is a brief description to let users know about the List.
	ListFieldDescription ListField = "description"
	// ListFieldOwnerID is unique identifier of this List's owner.
	ListFieldOwnerID ListField = "owner_id"
)

func listFieldStringArray(arr []ListField) []string {
	strs := make([]string, len(arr))
	for i, field := range arr {
		strs[i] = string(field)
	}
	return strs
}

// ListObj is the metadata for twitter lists
type ListObj struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	CreatedAt     string `json:"created_at"`
	Description   string `json:"description"`
	FollowerCount int    `json:"follower_count"`
	MemberCount   int    `json:"member_count"`
	Private       bool   `json:"private"`
	OwnerID       string `json:"owner_id"`
}

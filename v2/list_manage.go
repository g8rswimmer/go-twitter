package twitter

// ListMetaData is a list meta data
type ListMetaData struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Private     *bool   `json:"private"`
}

// ListCreateResponse is the response to creating a list
type ListCreateResponse struct {
	List *ListCreateData `json:"data"`
}

// ListCreateData is the data returned from creating a list
type ListCreateData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ListUpdateResponse is the response to updating a list
type ListUpdateResponse struct {
	List *ListUpdateData `json:"data"`
}

// ListUpdateData is the data returned from updating a list
type ListUpdateData struct {
	Updated bool `json:"updated"`
}

// ListDeleteResponse is the response to deleting a list
type ListDeleteResponse struct {
	List *ListDeleteData `json:"data"`
}

// ListDeleteData is the data returned from deleting a list
type ListDeleteData struct {
	Deleted bool `json:"deleted"`
}

// ListMemberData is the list member data
type ListMemberData struct {
	Member bool `json:"is_member"`
}

// ListAddMemberResponse is the list add member response
type ListAddMemberResponse struct {
	List *ListMemberData `json:"data"`
}

// ListRemoveMemberResponse is the list remove member response
type ListRemoveMemberResponse struct {
	List *ListMemberData `json:"data"`
}

// UserPinListData pinned data
type UserPinListData struct {
	Pinned bool `json:"pinned"`
}

// UserPinListResponse pin list response
type UserPinListResponse struct {
	List *UserPinListData `json:"data"`
}

// UserUnpinListResponse upin list response
type UserUnpinListResponse struct {
	List *UserPinListData `json:"data"`
}

type UserFollowListData struct {
	Following bool `json:"following"`
}

type UserFollowListResponse struct {
	List *UserFollowListData `json:"data"`
}

type UserUnfollowListResponse struct {
	List *UserFollowListData `json:"data"`
}

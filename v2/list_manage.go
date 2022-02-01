package twitter

type ListManageRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Private     *bool   `json:"private"`
}

type ListCreateResponse struct {
	List *ListCreateData `json:"data"`
}

type ListCreateData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ListUpdateResponse struct {
	List *ListUpdateData `json:"data"`
}

type ListUpdateData struct {
	Updated bool `json:"updated"`
}

type ListDeleteResponse struct {
	List *ListDeleteData `json:"data"`
}

type ListDeleteData struct {
	Deleted bool `json:"deleted"`
}

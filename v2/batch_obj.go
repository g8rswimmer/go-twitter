package twitter 


type BatchObj struct {
	Resumable bool `json:"resumable"`
	Type string `json:"type"`
	ID string `json:"id"`
	CreatedAt string `json:"created_at"`
	Name string `json:"name"`
	UploadURL string `json:"upload_url"`
	UploadExpiresAt string `json:"upload_expires_at"`
	DownloadURL string `json:"download_url"`
	DownloadExpiresAt string `json:"download_expires_at"`
	Status string `json:"status"`
	Error string `json:"error"`
}
package twitter

type ComplianceBatchJobStatus string

const (
	ComplianceBatchJobStatusCreated    ComplianceBatchJobStatus = "created"
	ComplianceBatchJobStatusComplete   ComplianceBatchJobStatus = "complete"
	ComplianceBatchJobStatusInProgress ComplianceBatchJobStatus = "in_progress"
	ComplianceBatchJobStatusFailed     ComplianceBatchJobStatus = "failed"
	ComplianceBatchJobStatusExpired    ComplianceBatchJobStatus = "expired"
)

type ComplianceBatchJobType string

const (
	ComplianceBatchJobTypeTweets ComplianceBatchJobType = "tweets"
	ComplianceBatchJobTypeUsers  ComplianceBatchJobType = "users"
)

type ComplianceBatchJobObj struct {
	Resumable         bool                     `json:"resumable"`
	Type              ComplianceBatchJobType   `json:"type"`
	ID                string                   `json:"id"`
	CreatedAt         string                   `json:"created_at"`
	Name              string                   `json:"name"`
	UploadURL         string                   `json:"upload_url"`
	UploadExpiresAt   string                   `json:"upload_expires_at"`
	DownloadURL       string                   `json:"download_url"`
	DownloadExpiresAt string                   `json:"download_expires_at"`
	Status            ComplianceBatchJobStatus `json:"status"`
	Error             string                   `json:"error"`
}

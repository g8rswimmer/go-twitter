package twitter

import "net/http"

// ComplianceBatchJobRaw is the raw data from a response
type ComplianceBatchJobRaw struct {
	Job *ComplianceBatchJobObj `json:"data"`
}

// CreateComplianceBatchJobOpts are the create compliance batch job options
type CreateComplianceBatchJobOpts struct {
	Name      string
	Resumable bool
}

// CreateComplianceBatchJobResponse is the response from creating a compliance batch job
type CreateComplianceBatchJobResponse struct {
	Raw       *ComplianceBatchJobRaw
	RateLimit *RateLimit
}

// ComplianceBatchJobResponse is the compliance batch job response
type ComplianceBatchJobResponse struct {
	Raw       *ComplianceBatchJobRaw
	RateLimit *RateLimit
}

// ComplianceBatchJobLookupOpts is the compliance batch lookups options
type ComplianceBatchJobLookupOpts struct {
	Status ComplianceBatchJobStatus
}

func (c ComplianceBatchJobLookupOpts) addQuery(req *http.Request) {
	q := req.URL.Query()
	if len(c.Status) > 0 {
		q.Add("status", string(c.Status))
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}

// ComplianceBatchJobLookupResponse is the response from the compliance batch lookup
type ComplianceBatchJobLookupResponse struct {
	Raw       *ComplianceBatchJobsRaw
	RateLimit *RateLimit
}

// ComplianceBatchJobsRaw is the jobs
type ComplianceBatchJobsRaw struct {
	Jobs []*ComplianceBatchJobObj `json:"data"`
}

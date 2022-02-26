package twitter

import "net/http"

type ComplianceBatchJobRaw struct {
	Job *ComplianceBatchJobObj `json:"data"`
}

type CreateComplianceBatchJobResponse struct {
	Raw       *ComplianceBatchJobRaw
	RateLimit *RateLimit
}

type ComplianceBatchJobResponse struct {
	Raw       *ComplianceBatchJobRaw
	RateLimit *RateLimit
}

type ComplianceBatchJobLookupOpts struct {
	Status ComplianceBatchJobStatus
}

func (c ComplianceBatchJobLookupOpts) addQuery(req *http.Request) {
	q := req.URL.Query()
	if len(c.Status) > 0 {
		q.Add("status", c.Status)
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}

type ComplianceBatchJobLookupResponse struct {
	Raw       *ComplianceBatchJobsRaw
	RateLimit *RateLimit
}

type ComplianceBatchJobsRaw struct {
	Jobs []*ComplianceBatchJobObj `json:"data"`
}

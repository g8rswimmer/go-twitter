package twitter

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ComplianceBatchJobStatus is the compliance batch job status
type ComplianceBatchJobStatus string

const (
	// ComplianceBatchJobStatusCreated is the created status
	ComplianceBatchJobStatusCreated ComplianceBatchJobStatus = "created"
	// ComplianceBatchJobStatusComplete is the complete status
	ComplianceBatchJobStatusComplete ComplianceBatchJobStatus = "complete"
	// ComplianceBatchJobStatusInProgress is the in_progress status
	ComplianceBatchJobStatusInProgress ComplianceBatchJobStatus = "in_progress"
	// ComplianceBatchJobStatusFailed is the failed status
	ComplianceBatchJobStatusFailed ComplianceBatchJobStatus = "failed"
	// ComplianceBatchJobStatusExpired is the expired status
	ComplianceBatchJobStatusExpired ComplianceBatchJobStatus = "expired"
)

// ComplianceBatchJobType is the compliance batch job type
type ComplianceBatchJobType string

const (
	// ComplianceBatchJobTypeTweets is the tweets job
	ComplianceBatchJobTypeTweets ComplianceBatchJobType = "tweets"
	// ComplianceBatchJobTypeUsers is the users job
	ComplianceBatchJobTypeUsers ComplianceBatchJobType = "users"
)

// ComplianceBatchJobResult is the downloaded result
type ComplianceBatchJobResult struct {
	ID         string `json:"id"`
	Action     string `json:"action"`
	CreatedAt  string `json:"created_at"`
	RedactedAt string `json:"redacted_at"`
	Reason     string `json:"reason"`
}

// ComplianceBatchJobDownloadResponse is the response from dowload results
type ComplianceBatchJobDownloadResponse struct {
	Results   []*ComplianceBatchJobResult
	RateLimit *RateLimit
}

// ComplianceBatchJobObj is the compliance batch job
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
	client            *http.Client
}

// Upload will upload ids from a reader
func (c ComplianceBatchJobObj) Upload(ctx context.Context, ids io.Reader) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.UploadURL, ids)
	if err != nil {
		return fmt.Errorf("compliance batch job upload request: %w", err)
	}
	req.Header.Add("Content-Type", "text/plain")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("compliance batch job upload response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return e
	}
	return nil
}

// Download will download the results of the job
func (c ComplianceBatchJobObj) Download(ctx context.Context) (*ComplianceBatchJobDownloadResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.DownloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("compliance batch job download request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("compliance batch job download response: %w", err)
	}
	defer resp.Body.Close()

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := json.NewDecoder(resp.Body).Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	results := []*ComplianceBatchJobResult{}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Split(batchResultsSeparator)
	for scanner.Scan() {
		result := &ComplianceBatchJobResult{}
		if err := json.Unmarshal(scanner.Bytes(), result); err != nil {
			return nil, &ResponseDecodeError{
				Name:      "compliance batch job download",
				Err:       err,
				RateLimit: rl,
			}
		}
		results = append(results, result)
	}

	return &ComplianceBatchJobDownloadResponse{
		Results:   results,
		RateLimit: rl,
	}, nil

}

func batchResultsSeparator(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if idx := bytes.Index(data, []byte("\r\n")); idx != -1 {
		return idx + len("\r\n"), data[0:idx], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

package twitter

import (
	"net/http"
	"time"
)

// Granularity is the granularity that you want the timeseries count data to be grouped by
type Granularity string

const (
	// GranularityMinute will group tweet in minutes
	GranularityMinute Granularity = "minute"
	// GranularityHour is the default granularity
	GranularityHour Granularity = "hour"
	// GranularityDay will group tweet on a daily basis
	GranularityDay Granularity = "day"
)

// TweetRecentCountsOpts are the optional paramters that can be passed to the tweet recent counts callout
type TweetRecentCountsOpts struct {
	StartTime   time.Time
	EndTime     time.Time
	SinceID     string
	UntilID     string
	Granularity Granularity
}

func (t TweetRecentCountsOpts) addQuery(req *http.Request) {
	q := req.URL.Query()
	if !t.StartTime.IsZero() {
		q.Add("start_time", t.StartTime.Format(time.RFC3339))
	}
	if !t.EndTime.IsZero() {
		q.Add("end_time", t.EndTime.Format(time.RFC3339))
	}
	if len(t.SinceID) > 0 {
		q.Add("since_id", t.SinceID)
	}
	if len(t.UntilID) > 0 {
		q.Add("until_id", t.UntilID)
	}
	if len(t.Granularity) > 0 {
		q.Add("granularity", string(t.Granularity))
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}

// TweetAllCountsOpts are the optional paramters that can be passed to the tweet all counts callout
type TweetAllCountsOpts struct {
	StartTime   time.Time
	EndTime     time.Time
	SinceID     string
	UntilID     string
	Granularity Granularity
	NextToken   string
}

func (t TweetAllCountsOpts) addQuery(req *http.Request) {
	q := req.URL.Query()
	if !t.StartTime.IsZero() {
		q.Add("start_time", t.StartTime.Format(time.RFC3339))
	}
	if !t.EndTime.IsZero() {
		q.Add("end_time", t.EndTime.Format(time.RFC3339))
	}
	if len(t.SinceID) > 0 {
		q.Add("since_id", t.SinceID)
	}
	if len(t.UntilID) > 0 {
		q.Add("until_id", t.UntilID)
	}
	if len(t.Granularity) > 0 {
		q.Add("granularity", string(t.Granularity))
	}
	if len(t.NextToken) > 0 {
		q.Add("next_token", t.NextToken)
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}

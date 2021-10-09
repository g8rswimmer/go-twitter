package twitter

import (
	"net/http"
	"time"
)

// TweetRecentCountsOpts are the optional paramters that can be passed to the tweet recent counts callout
type TweetRecentCountsOpts struct {
	StartTime   time.Time
	EndTime     time.Time
	SinceID     string
	UntilID     string
	Granularity string
}

func (t TweetRecentCountsOpts) addQuery(req *http.Request) {
	q := req.URL.Query()
	if t.StartTime.IsZero() == false {
		q.Add("start_time", t.StartTime.Format(time.RFC3339))
	}
	if t.EndTime.IsZero() == false {
		q.Add("end_time", t.EndTime.Format(time.RFC3339))
	}
	if len(t.SinceID) > 0 {
		q.Add("since_id", t.SinceID)
	}
	if len(t.UntilID) > 0 {
		q.Add("until_id", t.UntilID)
	}
	if len(t.Granularity) > 0 {
		q.Add("granularity", t.Granularity)
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}

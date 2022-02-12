package twitter

import (
	"net/http"
	"strconv"
	"time"
)

const (
	rateLimit     = "x-rate-limit-limit"
	rateRemaining = "x-rate-limit-remaining"
	rateReset     = "x-rate-limit-reset"
)

type Epoch int

func (e Epoch) Time() time.Time {
	return time.Unix(int64(e), 0)
}

type RateLimit struct {
	Limit     int
	Remaining int
	Reset     Epoch
}

func rateFromHeader(header http.Header) *RateLimit {
	limit, err := strconv.Atoi(header.Get(rateLimit))
	if err != nil {
		return nil
	}
	remaining, err := strconv.Atoi(header.Get(rateRemaining))
	if err != nil {
		return nil
	}
	reset, err := strconv.Atoi(header.Get(rateReset))
	if err != nil {
		return nil
	}
	return &RateLimit{
		Limit:     limit,
		Remaining: remaining,
		Reset:     Epoch(reset),
	}
}

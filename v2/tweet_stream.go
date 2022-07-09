package twitter

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// SystemMessageType stream system message types
type SystemMessageType string

// StreamErrorType is the type of streaming error
type StreamErrorType string

const (
	// InfoMessageType is the information system message type
	InfoMessageType SystemMessageType = "info"
	// WarnMessageType is the warning system message type
	WarnMessageType SystemMessageType = "warn"
	// ErrorMessageType is the error system message type
	ErrorMessageType SystemMessageType = "error"

	tweetStart  = "data"
	keepAliveTO = 21 * time.Second

	// TweetErrorType represents the tweet stream errors
	TweetErrorType StreamErrorType = "tweet"
	// SystemErrorType represents the system stream errors
	SystemErrorType StreamErrorType = "system"
)

// TweetSampleStreamOpts are the options for sample tweet stream
type TweetSampleStreamOpts struct {
	BackfillMinutes int
	Expansions      []Expansion
	MediaFields     []MediaField
	PlaceFields     []PlaceField
	PollFields      []PollField
	TweetFields     []TweetField
	UserFields      []UserField
}

func (t TweetSampleStreamOpts) addQuery(req *http.Request) {
	q := req.URL.Query()
	if len(t.Expansions) > 0 {
		q.Add("expansions", strings.Join(expansionStringArray(t.Expansions), ","))
	}
	if len(t.MediaFields) > 0 {
		q.Add("media.fields", strings.Join(mediaFieldStringArray(t.MediaFields), ","))
	}
	if len(t.PlaceFields) > 0 {
		q.Add("place.fields", strings.Join(placeFieldStringArray(t.PlaceFields), ","))
	}
	if len(t.PollFields) > 0 {
		q.Add("poll.fields", strings.Join(pollFieldStringArray(t.PollFields), ","))
	}
	if len(t.TweetFields) > 0 {
		q.Add("tweet.fields", strings.Join(tweetFieldStringArray(t.TweetFields), ","))
	}
	if len(t.UserFields) > 0 {
		q.Add("user.fields", strings.Join(userFieldStringArray(t.UserFields), ","))
	}
	if t.BackfillMinutes > 0 {
		q.Add("backfill_minutes", strconv.Itoa(t.BackfillMinutes))
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}

// TweetSearchStreamOpts are the options for the search stream
type TweetSearchStreamOpts struct {
	BackfillMinutes int
	Expansions      []Expansion
	MediaFields     []MediaField
	PlaceFields     []PlaceField
	PollFields      []PollField
	TweetFields     []TweetField
	UserFields      []UserField
}

func (t TweetSearchStreamOpts) addQuery(req *http.Request) {
	q := req.URL.Query()
	if len(t.Expansions) > 0 {
		q.Add("expansions", strings.Join(expansionStringArray(t.Expansions), ","))
	}
	if len(t.MediaFields) > 0 {
		q.Add("media.fields", strings.Join(mediaFieldStringArray(t.MediaFields), ","))
	}
	if len(t.PlaceFields) > 0 {
		q.Add("place.fields", strings.Join(placeFieldStringArray(t.PlaceFields), ","))
	}
	if len(t.PollFields) > 0 {
		q.Add("poll.fields", strings.Join(pollFieldStringArray(t.PollFields), ","))
	}
	if len(t.TweetFields) > 0 {
		q.Add("tweet.fields", strings.Join(tweetFieldStringArray(t.TweetFields), ","))
	}
	if len(t.UserFields) > 0 {
		q.Add("user.fields", strings.Join(userFieldStringArray(t.UserFields), ","))
	}
	if t.BackfillMinutes > 0 {
		q.Add("backfill_minutes", strconv.Itoa(t.BackfillMinutes))
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}

// StreamError is the error from the streaming
type StreamError struct {
	Type StreamErrorType
	Msg  string
	Err  error
}

func (e StreamError) Error() string {
	msg := fmt.Sprintf("%s: %s", e.Type, e.Msg)
	if e.Err == nil {
		return msg
	}
	return fmt.Sprintf("%s %s", msg, e.Err.Error())
}

// Is will compare the error against the stream error and type
func (e *StreamError) Is(target error) bool {
	cmp, ok := target.(*StreamError)
	if !ok {
		return false
	}
	return cmp.Type == e.Type
}

// Unwrap will return any error associated
func (e *StreamError) Unwrap() error {
	return e.Err
}

// TweetMessage is the tweet stream message
type TweetMessage struct {
	Raw *TweetRaw
}

// SystemMessage is the system stream message
type SystemMessage struct {
	Message string    `json:"message"`
	Sent    time.Time `json:"sent"`
}

// TweetStream is the stream handler
type TweetStream struct {
	tweets    chan *TweetMessage
	system    chan map[SystemMessageType]SystemMessage
	close     chan bool
	err       chan error
	alive     bool
	mutex     sync.RWMutex
	RateLimit *RateLimit
}

// StartTweetStream will start the tweet streaming
func StartTweetStream(stream io.ReadCloser) *TweetStream {
	ts := &TweetStream{
		tweets: make(chan *TweetMessage, 10),
		system: make(chan map[SystemMessageType]SystemMessage, 10),
		close:  make(chan bool),
		err:    make(chan error),
		mutex:  sync.RWMutex{},
		alive:  true,
	}

	go ts.handle(stream)

	return ts
}

func (ts *TweetStream) heartbeat(beat bool) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	ts.alive = beat
}

// Connection returns if the connect is still alive
func (ts *TweetStream) Connection() bool {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()
	return ts.alive
}

func (ts *TweetStream) handle(stream io.ReadCloser) {
	defer stream.Close()
	defer close(ts.tweets)
	defer close(ts.system)
	defer close(ts.close)
	defer close(ts.err)

	scanner := bufio.NewScanner(stream)
	scanner.Split(streamSeparator)
	timer := time.NewTimer(keepAliveTO)
	for {
		select {
		case <-ts.close:
			return
		case <-timer.C:
			ts.heartbeat(false)
		default:
		}

		if !scanner.Scan() {
			continue
		}

		timer.Stop()
		timer.Reset(keepAliveTO)
		ts.heartbeat(true)

		msg := scanner.Bytes()

		if len(msg) == 0 {
			continue
		}

		msgMap := map[string]interface{}{}
		if err := json.Unmarshal(msg, &msgMap); err != nil {
			select {
			case ts.err <- fmt.Errorf("stream error: unmarshal error %w", err):
			default:
			}
			continue
		}

		if _, tweet := msgMap[tweetStart]; tweet {
			single := &tweetraw{}
			if err := json.Unmarshal(msg, single); err != nil {
				sErr := &StreamError{
					Type: TweetErrorType,
					Msg:  "unmarshal tweet stream",
					Err:  err,
				}
				select {
				case ts.err <- sErr:
				default:
				}
				continue
			}
			raw := &TweetRaw{}
			raw.Tweets = make([]*TweetObj, 1)
			raw.Tweets[0] = single.Tweet
			raw.Includes = single.Includes
			raw.Errors = single.Errors

			tweetMsg := &TweetMessage{
				Raw: raw,
			}

			select {
			case ts.tweets <- tweetMsg:
			default:
			}
			continue
		}

		sysMsg := map[SystemMessageType]SystemMessage{}
		if err := json.Unmarshal(msg, &sysMsg); err != nil {
			sErr := &StreamError{
				Type: SystemErrorType,
				Msg:  "unmarshal system stream",
				Err:  err,
			}
			select {
			case ts.err <- sErr:
			default:
			}
			continue
		}
		select {
		case ts.system <- sysMsg:
		default:
		}
	}
}

// Tweets will return the channel to receive tweet stream messages
func (ts *TweetStream) Tweets() <-chan *TweetMessage {
	return ts.tweets
}

// SystemMessages will return the channel to receive system stream messages
func (ts *TweetStream) SystemMessages() <-chan map[SystemMessageType]SystemMessage {
	return ts.system
}

// Err will return the channel to receive any stream errors
func (ts *TweetStream) Err() <-chan error {
	return ts.err
}

// Close will close the stream and all channels
func (ts *TweetStream) Close() {
	ts.close <- true
}

func streamSeparator(data []byte, atEOF bool) (int, []byte, error) {
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

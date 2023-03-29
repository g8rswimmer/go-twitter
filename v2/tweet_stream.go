package twitter

import (
	"bufio"
	"bytes"
	"encoding/base64"
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

type streamType int

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
	// DisconnectErrorType represents the disconnection errors
	DisconnectErrorType StreamErrorType = "disconnect"

	disconnectionErrorsKey = "errors"
	disconnectionTitleKey  = "title"

	decodeErrStream   streamType = -1
	tweetStream       streamType = 1
	systemMsgStream   streamType = 2
	disconnectionErrs streamType = 3
	disconnectionErr  streamType = 4
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

// DisconnectionError contains the disconnection messages
type DisconnectionError struct {
	Disconnections []*Disconnection
	Connections    []*Connection
}

// Disconnection has the disconnection error
type Disconnection struct {
	Title          string `json:"title"`
	DisconnectType string `json:"disconnect_type"`
	Detail         string `json:"detail"`
	Type           string `json:"type"`
}

// Connection has the connection error
type Connection struct {
	Title           string `json:"title"`
	ConnectionIssue string `json:"connection_issue"`
	Detail          string `json:"detail"`
	Type            string `json:"type"`
}

type disconnection struct {
	Title           string `json:"title"`
	DisconnectType  string `json:"disconnect_type"`
	ConnectionIssue string `json:"connection_issue"`
	Detail          string `json:"detail"`
	Type            string `json:"type"`
}

func (d disconnection) disconnectType() bool {
	return d.DisconnectType != ""
}

func (d disconnection) toDisconnection() *Disconnection {
	return &Disconnection{
		Title:          d.Title,
		DisconnectType: d.DisconnectType,
		Detail:         d.Detail,
		Type:           d.Type,
	}
}

func (d disconnection) toConnection() *Connection {
	return &Connection{
		Title:           d.Title,
		ConnectionIssue: d.ConnectionIssue,
		Detail:          d.Detail,
		Type:            d.Type,
	}
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
	tweets        chan *TweetMessage
	system        chan map[SystemMessageType]SystemMessage
	disconnection chan *DisconnectionError
	close         chan bool
	err           chan error
	alive         bool
	mutex         sync.RWMutex
	RateLimit     *RateLimit
}

// StartTweetStream will start the tweet streaming
func StartTweetStream(stream io.ReadCloser) *TweetStream {
	ts := &TweetStream{
		tweets:        make(chan *TweetMessage, 10),
		system:        make(chan map[SystemMessageType]SystemMessage, 10),
		disconnection: make(chan *DisconnectionError, 10),
		close:         make(chan bool),
		err:           make(chan error),
		mutex:         sync.RWMutex{},
		alive:         true,
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

		reader, err := normalizeStream(msg)
		if err != nil {
			select {
			case ts.err <- fmt.Errorf("stream error: normalize error %w", err):
			default:
			}
			continue
		}

		sType, err := decodeStreamType(reader)
		if err != nil {
			select {
			case ts.err <- fmt.Errorf("stream error: unmarshal error %w", err):
			default:
			}
			continue
		}
		if _, err := reader.Seek(0, io.SeekStart); err != nil {
			select {
			case ts.err <- fmt.Errorf("stream error: seek error %w", err):
			default:
			}
			continue
		}
		decoder := json.NewDecoder(reader)

		switch sType {
		case tweetStream:
			ts.handleTweet(decoder)
		case systemMsgStream:
			ts.handleSystemMessage(decoder)
		case disconnectionErrs:
			ts.handleDisconnectErrors(decoder)
		case disconnectionErr:
			ts.handleDisconnectError(decoder)
		default:
		}
	}
}

func (ts *TweetStream) handleTweet(decoder *json.Decoder) {
	single := &tweetraw{}
	if err := decoder.Decode(single); err != nil {
		sErr := &StreamError{
			Type: TweetErrorType,
			Msg:  "unmarshal tweet stream",
			Err:  err,
		}
		select {
		case ts.err <- sErr:
		default:
		}
		return
	}
	raw := &TweetRaw{}
	raw.Tweets = make([]*TweetObj, 1)
	raw.Tweets[0] = single.Tweet
	raw.Includes = single.Includes
	raw.Errors = single.Errors
	raw.MatchingRules = single.MatchingRules

	tweetMsg := &TweetMessage{
		Raw: raw,
	}

	select {
	case ts.tweets <- tweetMsg:
	default:
	}
}

func (ts *TweetStream) handleSystemMessage(decoder *json.Decoder) {
	sysMsg := map[SystemMessageType]SystemMessage{}
	if err := decoder.Decode(&sysMsg); err != nil {
		sErr := &StreamError{
			Type: SystemErrorType,
			Msg:  "unmarshal system stream",
			Err:  err,
		}
		select {
		case ts.err <- sErr:
		default:
		}
		return
	}
	select {
	case ts.system <- sysMsg:
	default:
	}

}

func (ts *TweetStream) handleDisconnectErrors(decoder *json.Decoder) {
	disErrs := struct {
		Errors []disconnection `json:"errors"`
	}{}
	if err := decoder.Decode(&disErrs); err != nil {
		sErr := &StreamError{
			Type: DisconnectErrorType,
			Msg:  "unmarshal disconnect stream",
			Err:  err,
		}
		select {
		case ts.err <- sErr:
		default:
		}
		return
	}

	ds := &DisconnectionError{
		Disconnections: []*Disconnection{},
		Connections:    []*Connection{},
	}
	for _, d := range disErrs.Errors {
		switch {
		case d.disconnectType():
			ds.Disconnections = append(ds.Disconnections, d.toDisconnection())
		default:
			ds.Connections = append(ds.Connections, d.toConnection())
		}
	}

	select {
	case ts.disconnection <- ds:
	default:
	}
}

func (ts *TweetStream) handleDisconnectError(decoder *json.Decoder) {
	d := disconnection{}
	if err := decoder.Decode(&d); err != nil {
		sErr := &StreamError{
			Type: DisconnectErrorType,
			Msg:  "unmarshal disconnect stream",
			Err:  err,
		}
		select {
		case ts.err <- sErr:
		default:
		}
		return
	}

	ds := &DisconnectionError{
		Disconnections: []*Disconnection{},
		Connections:    []*Connection{},
	}
	switch {
	case d.disconnectType():
		ds.Disconnections = append(ds.Disconnections, d.toDisconnection())
	default:
		ds.Connections = append(ds.Connections, d.toConnection())
	}

	select {
	case ts.disconnection <- ds:
	default:
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

// DisconnectionError will return the channel to receive disconnect error messages
func (ts *TweetStream) DisconnectionError() <-chan *DisconnectionError {
	return ts.disconnection
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

func decodeStreamType(reader io.Reader) (streamType, error) {
	mm := map[string]interface{}{}
	if err := json.NewDecoder(reader).Decode(&mm); err != nil {
		return decodeErrStream, fmt.Errorf("decode stream type: %w", err)
	}
	for k := range mm {
		switch k {
		case tweetStart:
			return tweetStream, nil
		case string(InfoMessageType), string(WarnMessageType), string(ErrorMessageType):
			return systemMsgStream, nil
		case disconnectionErrorsKey:
			return disconnectionErrs, nil
		case disconnectionTitleKey:
			return disconnectionErr, nil
		default:
		}
	}
	return decodeErrStream, fmt.Errorf("decode stream message")
}

func normalizeStream(msg []byte) (*bytes.Reader, error) {
	str := string(msg)
	switch {
	case strings.Contains(str, ":"):
		return bytes.NewReader(msg), nil
	default:
		decodedMsg, err := base64.StdEncoding.DecodeString(str)
		if err != nil {
			return nil, fmt.Errorf("stream normalize stream base64: %w", err)
		}
		return bytes.NewReader(decodedMsg), nil
	}
}

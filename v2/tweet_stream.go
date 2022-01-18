package twitter

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

	tweetStart = "data"

	// TweetErrorType represents the tweet stream errrors
	TweetErrorType StreamErrorType = "tweet"
	// SystemErrorType represents the system stream errors
	SystemErrorType StreamErrorType = "system"
)

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

// Unwrap will return any error assocaited
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
	tweets chan *TweetMessage
	system chan map[SystemMessageType]SystemMessage
	close  chan bool
	err    chan error
}

// StartTweetStream will start the tweet streaming
func StartTweetStream(stream io.ReadCloser) *TweetStream {
	ts := &TweetStream{
		tweets: make(chan *TweetMessage, 10),
		system: make(chan map[SystemMessageType]SystemMessage, 10),
		close:  make(chan bool),
		err:    make(chan error),
	}

	go ts.handle(stream)

	return ts
}

func (ts *TweetStream) handle(stream io.ReadCloser) {
	defer stream.Close()
	defer close(ts.tweets)
	defer close(ts.system)
	defer close(ts.close)
	defer close(ts.err)

	scanner := bufio.NewScanner(stream)
	scanner.Split(streamSeperator)
	for {
		select {
		case <-ts.close:
			return
		default:
		}

		if !scanner.Scan() {
			continue
		}

		msg := scanner.Bytes()

		if len(msg) == 0 {
			continue
		}

		msgMap := map[string]interface{}{}
		if err := json.Unmarshal(msg, &msgMap); err != nil {
			ts.err <- fmt.Errorf("stream error: unmarshal error %w", err)
			continue
		}

		if _, tweet := msgMap[tweetStart]; tweet {
			single := &tweetraw{}
			if err := json.Unmarshal(msg, single); err != nil {
				ts.err <- &StreamError{
					Type: TweetErrorType,
					Msg:  "umarshal tweet stream",
					Err:  err,
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
			ts.tweets <- tweetMsg
			continue
		}

		sysMsg := map[SystemMessageType]SystemMessage{}
		if err := json.Unmarshal(msg, &sysMsg); err != nil {
			ts.err <- &StreamError{
				Type: SystemErrorType,
				Msg:  "umarshal system stream",
				Err:  err,
			}
			continue
		}
		ts.system <- sysMsg
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

func streamSeperator(data []byte, atEOF bool) (int, []byte, error) {
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

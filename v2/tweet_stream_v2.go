package twitter

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// StreamedTweet is the tweet stream message
type StreamedTweet struct {
	Tweet         *TweetObj         `json:"data,omitempty"`
	Includes      *TweetRawIncludes `json:"includes,omitempty"`
	Errors        []*ErrorObj       `json:"errors,omitempty"`
	MatchingRules []StreamRuleTag   `json:"matching_rules,omitempty"`
}

// StreamRuleTag describes which search rule the tweet was streamed for.
type StreamRuleTag struct {
	ID  string `json:"id"`
	Tag string `json:"tag"`
}

// TweetStreamV2 is the stream handler.
//
// Compared to TweetStream:
//   - Supports search rule tags
//   - Callback-based interface, instead of channels, so it's harder to mess up
//     the event loop by accident
//   - Does not silently drop events if your code is not keeping up with the volume
//   - Stops when the context is cancelled
//   - You control in which goroutine the event loop and your callbacks run
//   - Clear indication of when the connection is broken and needs to be restarted
type TweetStreamV2 struct {
	stream    io.ReadCloser
	RateLimit *RateLimit
}

// StartTweetStream creates a tweet stream object. The user must call Run() method
// order to start the event loop and actually receive tweets.
func StartTweetStreamV2(stream io.ReadCloser) *TweetStreamV2 {
	ts := &TweetStreamV2{stream: stream}

	return ts
}

// TweetStreamV2Options contains the options for running a tweet stream.
type TweetStreamV2Options struct {
	// OnTweet is a callback that gets called on every incoming tweet.
	OnTweet func(*StreamedTweet)
	// OnSystemMessage is a callback that gets called on every in-band system message.
	OnSystemMessage func(kind SystemMessageType, msg *SystemMessage)
	// OnTransientError is a callback that gets called on recoverable errors that
	// don't require restarting the stream.
	OnTransientError func(error)
}

func (opts TweetStreamV2Options) tweet(t *StreamedTweet) {
	if opts.OnTweet != nil {
		opts.OnTweet(t)
	}
}

func (opts TweetStreamV2Options) systemMessage(kind SystemMessageType, msg *SystemMessage) {
	if opts.OnSystemMessage != nil {
		opts.OnSystemMessage(kind, msg)
	}
}

func (opts TweetStreamV2Options) transientError(err error) {
	if opts.OnTransientError != nil {
		opts.OnTransientError(err)
	}
}

// Run runs the event loop of the stream. Returns when the stream is closed or failed.
func (ts *TweetStreamV2) Run(ctx context.Context, opts TweetStreamV2Options) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, 1)
	timer := time.NewTimer(keepAliveTO)
	defer timer.Stop()
	go func() {
		// Since scanner.Scan() is blocking, we need to watch the context and
		// heartbeat timer in a separate goroutine, and close the stream to
		// have scanner.Scan() return.
		defer ts.stream.Close()

		select {
		case <-ctx.Done():
			errCh <- ctx.Err()
		case <-timer.C:
			errCh <- fmt.Errorf("heartbeat timeout")
		}
	}()

	scanner := bufio.NewScanner(ts.stream)
	scanner.Split(streamSeparator)
	for {
		if !scanner.Scan() {
			select {
			case err := <-errCh:
				return err
			default:
				if err := scanner.Err(); err != nil {
					return fmt.Errorf("stream error: %w", err)
				}
				return fmt.Errorf("stream error: EOF")
			}
		}

		timer.Stop()
		timer.Reset(keepAliveTO)

		msg := scanner.Bytes()

		if len(msg) == 0 {
			continue
		}

		reader, err := normalizeStream(msg)
		if err != nil {
			opts.transientError(fmt.Errorf("stream error: normalize error %w", err))
		}

		sType, err := decodeStreamType(reader)
		if err != nil {
			opts.transientError(fmt.Errorf("stream error: unmarshal error %w", err))
		}
		if _, err := reader.Seek(0, io.SeekStart); err != nil {
			opts.transientError(fmt.Errorf("stream error: seek error %w", err))
		}
		decoder := json.NewDecoder(reader)

		switch sType {
		case tweetStream:
			ts.handleTweet(decoder, opts)
		case systemMsgStream:
			ts.handleSystemMessage(decoder, opts)
		case disconnectionErrs:
			return ts.handleDisconnectErrors(decoder)
		case disconnectionErr:
			return ts.handleDisconnectError(decoder)
		default:
			opts.transientError(fmt.Errorf("unsupported payload: %q", string(msg)))
		}
	}
}

func (ts *TweetStreamV2) handleTweet(decoder *json.Decoder, opts TweetStreamV2Options) {
	tweet := &StreamedTweet{}
	if err := decoder.Decode(tweet); err != nil {
		opts.transientError(&StreamError{
			Type: TweetErrorType,
			Msg:  "unmarshal tweet stream",
			Err:  err,
		})
		return
	}

	opts.tweet(tweet)
}

func (ts *TweetStreamV2) handleSystemMessage(decoder *json.Decoder, opts TweetStreamV2Options) {
	sysMsg := map[SystemMessageType]SystemMessage{}
	if err := decoder.Decode(&sysMsg); err != nil {
		opts.transientError(&StreamError{
			Type: SystemErrorType,
			Msg:  "unmarshal system stream",
			Err:  err,
		})
		return
	}
	for k, v := range sysMsg {
		v := v // copy loop variable
		opts.systemMessage(k, &v)
	}
}

func (ts *TweetStreamV2) handleDisconnectErrors(decoder *json.Decoder) error {
	disErrs := struct {
		Errors []disconnection `json:"errors"`
	}{}
	if err := decoder.Decode(&disErrs); err != nil {
		return &StreamError{
			Type: DisconnectErrorType,
			Msg:  "unmarshal disconnect stream",
			Err:  err,
		}
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

	return ds
}

func (ts *TweetStreamV2) handleDisconnectError(decoder *json.Decoder) error {
	d := disconnection{}
	if err := decoder.Decode(&d); err != nil {
		return &StreamError{
			Type: DisconnectErrorType,
			Msg:  "unmarshal disconnect stream",
			Err:  err,
		}
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

	return ds
}

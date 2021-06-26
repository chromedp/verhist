package omahaproxy

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/kenshaw/httplog"
)

// DefaultTransport is the default transport.
var DefaultTransport = http.DefaultTransport

// Client is a omaha proxy client.
type Client struct {
	Transport http.RoundTripper
}

// New creates a new omaha proxy client.
func New(opts ...Option) *Client {
	cl := &Client{
		Transport: DefaultTransport,
	}
	for _, o := range opts {
		o(cl)
	}
	return cl
}

// get retrieves data from the url.
func (cl *Client) get(ctx context.Context, urlstr string) ([]byte, error) {
	req, err := http.NewRequest("GET", urlstr, nil)
	if err != nil {
		return nil, err
	}
	// retrieve and decode
	httpClient := &http.Client{Transport: cl.Transport}
	res, err := httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not retrieve %s (status: %d)", urlstr, res.StatusCode)
	}
	return ioutil.ReadAll(res.Body)
}

// Recent retrieves the recent release history from the omaha proxy.
func (cl *Client) Recent(ctx context.Context) ([]Release, error) {
	buf, err := cl.get(ctx, "https://omahaproxy.appspot.com/history")
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(bytes.NewReader(buf))
	r.FieldsPerRecord = 4
	r.TrimLeadingSpace = true
	var history []Release
	var i int
loop:
	for {
		row, err := r.Read()
		switch {
		case err != nil && err == io.EOF:
			break loop
		case err != nil:
			return nil, err
		}
		i++
		if i == 1 {
			continue
		}
		timestamp, err := time.Parse("2006-01-02 15:04:05.999999999", row[3])
		if err != nil {
			return nil, err
		}
		history = append(history, Release{
			OS:        row[0],
			Channel:   row[1],
			Version:   row[2],
			Timestamp: timestamp,
		})
	}
	return history, nil
}

// Entries retrieves latest version entries from the omaha proxy.
func (cl *Client) Entries(ctx context.Context) ([]VersionEntry, error) {
	buf, err := cl.get(ctx, "https://omahaproxy.appspot.com/json")
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewReader(buf))
	dec.DisallowUnknownFields()
	var entries []VersionEntry
	if err := dec.Decode(&entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// Latest returns the latest version for the provided os and channel from the
// omaha proxy.
func (cl *Client) Latest(ctx context.Context, os, channel string) (Version, error) {
	entries, err := cl.Entries(ctx)
	if err != nil {
		return Version{}, err
	}
	for _, entry := range entries {
		if entry.OS != os {
			continue
		}
		for _, v := range entry.Versions {
			if v.Channel == channel {
				return v, nil
			}
		}
	}
	return Version{}, fmt.Errorf("could not find latest version for channel %s (%s)", channel, os)
}

// Release holds browser release information.
type Release struct {
	OS        string
	Channel   string
	Version   string
	Timestamp time.Time
}

// Version wraps browser version information.
type Version struct {
	BranchCommit       string `json:"branch_commit"`
	BranchBasePosition string `json:"branch_base_position"`
	SkiaCommit         string `json:"skia_commit"`
	V8Version          string `json:"v8_version"`
	PreviousVersion    string `json:"previous_version"`
	V8Commit           string `json:"v8_commit"`
	TrueBranch         string `json:"true_branch"`
	PreviousReldate    string `json:"previous_reldate"`
	BranchBaseCommit   string `json:"branch_base_commit"`
	Version            string `json:"version"`
	CurrentReldate     string `json:"current_reldate"`
	CurrentVersion     string `json:"current_version"`
	OS                 string `json:"os"`
	Channel            string `json:"channel"`
	ChromiumCommit     string `json:"chromium_commit"`
}

// String satisfies the fmt.Stringer interface.
func (v Version) String() string {
	return fmt.Sprintf("Chromium %s (v8: %s, os: %s, channel: %s)", v.Version, v.V8Version, v.OS, v.Channel)
}

// VersionEntry is a OS version entry detailing the available browser
// version entries.
type VersionEntry struct {
	OS       string    `json:"os"`
	Versions []Version `json:"versions"`
}

// Option is a omaha proxy client option.
type Option func(*Client)

// WithTransport is a omaha proxy client option to set the http transport.
func WithTransport(transport http.RoundTripper) Option {
	return func(cl *Client) {
		cl.Transport = transport
	}
}

// WithLogf is a transmission rpc client option to set a log handler for HTTP
// requests and responses.
func WithLogf(logf interface{}, opts ...httplog.Option) Option {
	return func(cl *Client) {
		cl.Transport = httplog.NewPrefixedRoundTripLogger(cl.Transport, logf, opts...)
	}
}

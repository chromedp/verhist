package verhist

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/kenshaw/httplog"
)

// https://versionhistory.googleapis.com/v1/chrome/platforms/all/channels/all/versions/
// https://versionhistory.googleapis.com/v1/chrome/platforms/all/channels/all/versions/all/releases?filter=endtime%3E2023-01-01T00:00:00Z

// https://developer.chrome.com/docs/web-platform/versionhistory/guide
// https://developer.chrome.com/docs/web-platform/versionhistory/reference
// https://developer.chrome.com/docs/web-platform/versionhistory/examples

// DefaultTransport is the default transport.
var DefaultTransport = http.DefaultTransport

// BaseURL is the base URL.
var BaseURL = "https://versionhistory.googleapis.com"

// Client is a version history client.
type Client struct {
	Transport http.RoundTripper
}

// New creates a new version history client.
func New(opts ...Option) *Client {
	cl := &Client{
		Transport: DefaultTransport,
	}
	for _, o := range opts {
		o(cl)
	}
	return cl
}

// Versions returns the versions for the os, channel.
func (cl *Client) Versions(ctx context.Context, os, channel string, q ...string) ([]Version, error) {
	if len(q) == 0 {
		q = []string{
			"order_by", "version desc",
		}
	}
	res := new(VersionsResponse)
	if err := grab(ctx, BaseURL+"/v1/chrome/platforms/"+os+"/channels/"+channel+"/versions", cl.Transport, res, q...); err != nil {
		return nil, err
	}
	return res.Versions, nil
}

// UserAgent builds the user agent for the os, channel.
func (cl *Client) UserAgent(ctx context.Context, os, channel string) (string, error) {
	versions, err := cl.Versions(ctx, os, channel)
	switch {
	case err != nil:
		return "", err
	case len(versions) == 0:
		return "", errors.New("no versions returned")
	}
	return versions[0].UserAgent(os), nil
}

// Option is a version history client option.
type Option func(*Client)

// WithTransport is a version history client option to set the http transport.
func WithTransport(transport http.RoundTripper) Option {
	return func(cl *Client) {
		cl.Transport = transport
	}
}

// WithLogf is a version history client option to set a log handler for HTTP
// requests and responses.
func WithLogf(logf interface{}, opts ...httplog.Option) Option {
	return func(cl *Client) {
		cl.Transport = httplog.NewPrefixedRoundTripLogger(cl.Transport, logf, opts...)
	}
}

// VersionsResponse wraps the versions API response.
type VersionsResponse struct {
	Versions      []Version `json:"versions,omitempty"`
	NextPageToken string    `json:"nextPageToken,omitempty"`
}

// Version contains information about a chrome release.
type Version struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// UserAgent builds the user agent for the
func (ver Version) UserAgent(os string) string {
	typ := "Windows NT 10.0; Win64; x64"
	switch strings.ToLower(os) {
	case "linux":
		typ = "X11; Linux x86_64"
	case "mac", "mac_arm64":
		typ = "Macintosh; Intel Mac OS X 10_15_7"
	}
	v := "120.0.0.0"
	if i := strings.Index(ver.Version, "."); i != -1 {
		v = ver.Version[:i] + ".0.0.0"
	}
	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", typ, v)
}

// grab grabs the url and json decodes it.
func grab(ctx context.Context, urlstr string, transport http.RoundTripper, v interface{}, q ...string) error {
	if len(q)%2 != 0 {
		return errors.New("invalid query")
	}
	z := make(url.Values)
	for i := 0; i < len(q); i += 2 {
		z.Add(q[i], q[i+1])
	}
	s := z.Encode()
	if s != "" {
		s = "?" + s
	}
	req, err := http.NewRequestWithContext(ctx, "GET", urlstr+s, nil)
	if err != nil {
		return err
	}
	cl := &http.Client{
		Transport: transport,
	}
	res, err := cl.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("could not retrieve %s (status: %d)", urlstr, res.StatusCode)
	}
	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

package verhist

import (
	"context"
	"encoding/json"
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

// Platforms returns platforms.
func (cl *Client) Platforms(ctx context.Context) ([]Platform, error) {
	res := new(PlatformsResponse)
	if err := grab(ctx, BaseURL+"/v1/chrome/platforms/", cl.Transport, res); err != nil {
		return nil, err
	}
	return res.Platforms, nil
}

// Channels returns channels for the platform type.
func (cl *Client) Channels(ctx context.Context, typ PlatformType) ([]Channel, error) {
	res := new(ChannelsResponse)
	if err := grab(ctx, BaseURL+"/v1/chrome/platforms/"+typ.String()+"/channels", cl.Transport, res); err != nil {
		return nil, err
	}
	return res.Channels, nil
}

// All returns all channels.
func (cl *Client) All(ctx context.Context) ([]Channel, error) {
	return cl.Channels(ctx, All)
}

// Versions returns the versions for the platform, channel.
func (cl *Client) Versions(ctx context.Context, platform, channel string, q ...string) ([]Version, error) {
	if len(q) == 0 {
		q = []string{
			"order_by", "version desc",
		}
	}
	res := new(VersionsResponse)
	if err := grab(ctx, BaseURL+"/v1/chrome/platforms/"+platform+"/channels/"+channel+"/versions", cl.Transport, res, q...); err != nil {
		return nil, err
	}
	return res.Versions, nil
}

// Latest veturns the latest version for the platform, channel.
func (cl *Client) Latest(ctx context.Context, platform, channel string) (Version, error) {
	versions, err := cl.Versions(ctx, platform, channel)
	switch {
	case err != nil:
		return Version{}, err
	case len(versions) == 0:
		return Version{}, ErrNoVersionsReturned
	}
	return versions[0], nil
}

// UserAgent builds the user agent for the platform, channel.
func (cl *Client) UserAgent(ctx context.Context, platform, channel string) (string, error) {
	latest, err := cl.Latest(ctx, platform, channel)
	if err != nil {
		return "", err
	}
	return latest.UserAgent(platform), nil
}

// PlatformsResponse wraps the platforms API response.
type PlatformsResponse struct {
	Platforms     []Platform `json:"platforms,omitempty"`
	NextPageToken string     `json:"nextPageToken,omitempty"`
}

// Platform contains information about a chrome platform.
type Platform struct {
	Name         string       `json:"name,omitempty"`
	PlatformType PlatformType `json:"platformType,omitempty"`
}

// PlatformType is a platform type.
type PlatformType string

// Platform types.
const (
	All          PlatformType = "all"
	Android      PlatformType = "android"
	ChromeOS     PlatformType = "chromeos"
	Fuchsia      PlatformType = "fuchsia"
	IOS          PlatformType = "ios"
	LacrosARM32  PlatformType = "lacros_arm32"
	LacrosARM64  PlatformType = "lacros_arm64"
	Lacros       PlatformType = "lacros"
	Linux        PlatformType = "linux"
	MacARM64     PlatformType = "mac_arm64"
	Mac          PlatformType = "mac"
	Webview      PlatformType = "webview"
	Windows64    PlatformType = "win64"
	WindowsARM64 PlatformType = "win_arm64"
	Windows      PlatformType = "win"
)

// String satisfies the [fmt.Stinger] interface.
func (typ PlatformType) String() string {
	return strings.ToLower(string(typ))
}

// MarshalText satisfies the [encoding.TextMarshaler] interface.
func (typ PlatformType) MarshalText() ([]byte, error) {
	return []byte(typ.String()), nil
}

// UnmarshalText satisfies the [encoding.TextUnmarshaler] interface.
func (typ *PlatformType) UnmarshalText(buf []byte) error {
	switch PlatformType(strings.ToLower(string(buf))) {
	case All:
		*typ = All
	case Android:
		*typ = Android
	case ChromeOS:
		*typ = ChromeOS
	case Fuchsia:
		*typ = Fuchsia
	case IOS:
		*typ = IOS
	case LacrosARM32:
		*typ = LacrosARM32
	case LacrosARM64:
		*typ = LacrosARM64
	case Lacros:
		*typ = Lacros
	case Linux:
		*typ = Linux
	case MacARM64:
		*typ = MacARM64
	case Mac:
		*typ = Mac
	case Webview:
		*typ = Webview
	case Windows64:
		*typ = Windows64
	case WindowsARM64:
		*typ = WindowsARM64
	case Windows:
		*typ = Windows
	default:
		return ErrInvalidPlatformType
	}
	return nil
}

// ChannelsResponse wraps the channels API response.
type ChannelsResponse struct {
	Channels      []Channel `json:"channels,omitempty"`
	NextPageToken string    `json:"nextPageToken,omitempty"`
}

// Channel contains information about a chrome channel.
type Channel struct {
	Name        string      `json:"name,omitempty"`
	ChannelType ChannelType `json:"channelType,omitempty"`
}

// ChannelType is a channel type.
type ChannelType string

// Channel types.
const (
	Beta       ChannelType = "beta"
	CanaryASAN ChannelType = "canary_asan"
	Canary     ChannelType = "canary"
	Dev        ChannelType = "dev"
	Extended   ChannelType = "extended"
	LTS        ChannelType = "lts"
	LTC        ChannelType = "ltc"
	Stable     ChannelType = "stable"
)

// String satisfies the [fmt.Stinger] interface.
func (typ ChannelType) String() string {
	return strings.ToLower(string(typ))
}

// MarshalText satisfies the [encoding.TextMarshaler] interface.
func (typ ChannelType) MarshalText() ([]byte, error) {
	return []byte(typ.String()), nil
}

// UnmarshalText satisfies the [encoding.TextUnmarshaler] interface.
func (typ *ChannelType) UnmarshalText(buf []byte) error {
	switch ChannelType(strings.ToLower(string(buf))) {
	case Beta:
		*typ = Beta
	case CanaryASAN:
		*typ = CanaryASAN
	case Canary:
		*typ = Canary
	case Dev:
		*typ = Dev
	case Extended:
		*typ = Extended
	case LTS:
		*typ = LTS
	case LTC:
		*typ = LTC
	case Stable:
		*typ = Stable
	default:
		return ErrInvalidChannelType
	}
	return nil
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

// UserAgent builds a user agent for the platform.
func (ver Version) UserAgent(platform string) string {
	typ, extra := "Windows NT 10.0; Win64; x64", ""
	switch strings.ToLower(platform) {
	case "linux":
		typ = "X11; Linux x86_64"
	case "mac", "mac_arm64":
		typ = "Macintosh; Intel Mac OS X 10_15_7"
	case "android":
		typ, extra = "Linux; Android 10; K", " Mobile"
	}
	v := "120.0.0.0"
	if i := strings.Index(ver.Version, "."); i != -1 {
		v = ver.Version[:i] + ".0.0.0"
	}
	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s%s Safari/537.36", typ, v, extra)
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

// grab grabs the url and json decodes it.
func grab(ctx context.Context, urlstr string, transport http.RoundTripper, v interface{}, q ...string) error {
	if len(q)%2 != 0 {
		return ErrInvalidQuery
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

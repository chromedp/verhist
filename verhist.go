// Package verhist provides a client and utilities for working with the
// Chrome version history API.
//
// See: https://developer.chrome.com/docs/web-platform/versionhistory/guide
package verhist

import (
	"context"
)

/*
// All returns all channels.
func All(ctx context.Context, opts ...Option) ([]Channel, error) {
	return New(opts...).All(ctx)
}
*/

// Platforms returns platforms.
func Platforms(ctx context.Context, opts ...Option) ([]Platform, error) {
	return New(opts...).Platforms(ctx)
}

// Channels returns channels for the platform type.
func Channels(ctx context.Context, typ PlatformType, opts ...Option) ([]Channel, error) {
	return New(opts...).Channels(ctx, typ)
}

// Versions returns the versions for the platform, channel.
func Versions(ctx context.Context, platform, channel string, opts ...Option) ([]Version, error) {
	return New(opts...).Versions(ctx, platform, channel)
}

// Latest veturns the latest version for the platform, channel.
func Latest(ctx context.Context, platform, channel string, opts ...Option) (Version, error) {
	return New(opts...).Latest(ctx, platform, channel)
}

// UserAgent builds the user agent for the platform, channel.
func UserAgent(ctx context.Context, platform, channel string, opts ...Option) (string, error) {
	return New(opts...).UserAgent(ctx, platform, channel)
}

// Error is a error.
type Error string

// Error satisfies the [error] interface.
func (err Error) Error() string {
	return string(err)
}

// Errors.
const (
	// ErrInvalidPlatformType is the invalid platform type error.
	ErrInvalidPlatformType Error = "invalid platform type"
	// ErrInvalidChannelType is the invalid channel type error.
	ErrInvalidChannelType Error = "invalid channel type"
	// ErrInvalidQuery is the invalid query error.
	ErrInvalidQuery Error = "invalid query"
	// ErrNoVersionsReturned is the no versions returned error.
	ErrNoVersionsReturned Error = "no versions returned"
)

// Package verhist provides a client and utilities for working with the
// Chrome version history API.
//
// See: https://developer.chrome.com/docs/web-platform/versionhistory/guide
package verhist

import (
	"context"
)

// Versions returns the versions for the os, channel.
func Versions(ctx context.Context, os, channel string, opts ...Option) ([]Version, error) {
	return New(opts...).Versions(ctx, os, channel)
}

// UserAgent builds the user agent for the os, channel.
func UserAgent(ctx context.Context, os, channel string, opts ...Option) (string, error) {
	return New(opts...).UserAgent(ctx, os, channel)
}

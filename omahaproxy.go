// Package omahaproxy provides a client and utilities for working with the
// Chrome Omaha Proxy.
//
// See: https://omahaproxy.appspot.com
package omahaproxy

import (
	"context"
)

// Recent returns the recent release information from the omaha proxy.
func Recent(ctx context.Context, opts ...Option) ([]Release, error) {
	return New(opts...).Recent(ctx)
}

// Entries returns the latest version entries from the omaha proxy.
func Entries(ctx context.Context, opts ...Option) ([]VersionEntry, error) {
	return New(opts...).Entries(ctx)
}

// Latest retrieves the latest version for the specified os and channel from
// the omaha proxy.
func Latest(ctx context.Context, os, channel string, opts ...Option) (Version, error) {
	return New(opts...).Latest(ctx, os, channel)
}

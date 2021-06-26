// Package omahaproxy provides a client and utilities for working with the
// Chrome Omaha Proxy.
//
// See: https://omahaproxy.appspot.com
package omahaproxy

import (
	"context"
)

// Recent returns the recent release information from the omaha proxy.
func Recent(opts ...Option) ([]Release, error) {
	return New(opts...).Recent(context.Background())
}

// Entries returns the latest version entries from the omaha proxy.
func Entries(opts ...Option) ([]VersionEntry, error) {
	return New(opts...).Entries(context.Background())
}

// Latest retrieves the latest version for the specified os and channel from
// the omaha proxy.
func Latest(os, channel string, opts ...Option) (Version, error) {
	return New(opts...).Latest(context.Background(), os, channel)
}

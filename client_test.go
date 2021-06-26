package omahaproxy

import (
	"context"
	"testing"
)

func TestRecent(t *testing.T) {
	t.Parallel()
	cl := New(WithLogf(t.Logf))
	releases, err := cl.Recent(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	for _, release := range releases {
		t.Logf("os: %s channel: %s version: %s\n", release.OS, release.Channel, release.Version)
	}
}

func TestLatest(t *testing.T) {
	t.Parallel()
	cl := New(WithLogf(t.Logf))
	ver, err := cl.Latest(context.Background(), "linux", "stable")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	t.Logf("latest: %v", ver)
}

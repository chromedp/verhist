package verhist

import (
	"context"
	"testing"
)

func TestUserAgent(t *testing.T) {
	t.Parallel()
	cl := New(WithLogf(t.Logf))
	userAgent, err := cl.UserAgent(context.Background(), "linux", "stable")
	switch {
	case err != nil:
		t.Fatalf("expected no error, got: %v", err)
	case userAgent == "":
		t.Errorf("expected non-empty user agent")
	}
	t.Logf("user agent: %v", userAgent)
}

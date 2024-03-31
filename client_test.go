package verhist

import (
	"context"
	"os"
	"testing"
)

func TestPlatforms(t *testing.T) {
	t.Parallel()
	var opts []Option
	if os.Getenv("VERBOSE") != "" {
		opts = append(opts, WithLogf(t.Logf))
	}
	platforms, err := Platforms(context.Background(), opts...)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	for _, platform := range platforms {
		t.Logf("name: %s platform: %s", platform.Name, platform.PlatformType)
	}
}

func TestChannels(t *testing.T) {
	t.Parallel()
	var opts []Option
	if os.Getenv("VERBOSE") != "" {
		opts = append(opts, WithLogf(t.Logf))
	}
	for _, typ := range platforms() {
		t.Run(typ.String(), func(t *testing.T) {
			t.Parallel()
			channels, err := Channels(context.Background(), typ, opts...)
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			for _, channel := range channels {
				t.Logf("name: %s channel: %s", channel.Name, channel.ChannelType)
			}
		})
	}
}

func TestVersions(t *testing.T) {
	t.Parallel()
	var opts []Option
	if os.Getenv("VERBOSE") != "" {
		opts = append(opts, WithLogf(t.Logf))
	}
	versions, err := Versions(context.Background(), "linux", "stable", opts...)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	for _, version := range versions {
		t.Logf("name: %s version: %s", version.Name, version.Version)
	}
}

func TestLatest(t *testing.T) {
	t.Parallel()
	var opts []Option
	if os.Getenv("VERBOSE") != "" {
		opts = append(opts, WithLogf(t.Logf))
	}
	latest, err := Latest(context.Background(), "win", "stable", opts...)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	t.Logf("name: %s version: %s", latest.Name, latest.Version)
}

func TestUserAgent(t *testing.T) {
	t.Parallel()
	var opts []Option
	if os.Getenv("VERBOSE") != "" {
		opts = append(opts, WithLogf(t.Logf))
	}
	userAgent, err := UserAgent(context.Background(), "linux", "stable", opts...)
	switch {
	case err != nil:
		t.Fatalf("expected no error, got: %v", err)
	case userAgent == "":
		t.Errorf("expected non-empty user agent")
	}
	t.Logf("user agent: %v", userAgent)
}

func platforms() []PlatformType {
	return []PlatformType{
		All,
		Android,
		ChromeOS,
		Fuchsia,
		IOS,
		LacrosARM32,
		LacrosARM64,
		Lacros,
		Linux,
		MacARM64,
		Mac,
		Webview,
		Windows64,
		WindowsARM64,
		Windows,
	}
}

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/chromedp/verhist"
)

func main() {
	latest := flag.Bool("latest", false, "latest")
	latestUA := flag.Bool("latest-ua", false, "latest ua")
	platform := flag.String("platform", "win64", "platform")
	channel := flag.String("channel", "stable", "channel")
	flag.Parse()
	if err := run(context.Background(), os.Stdout, *latest, *latestUA, *platform, *channel); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, w io.Writer, latest, latestUA bool, platform, channel string) error {
	if latestUA {
		userAgent, err := verhist.UserAgent(ctx, platform, channel)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(w, userAgent)
		return err
	}
	versions, err := verhist.Versions(ctx, platform, channel)
	switch {
	case err != nil:
		return err
	case len(versions) == 0:
		return verhist.ErrNoVersionsAvailable
	case latest:
		_, err := fmt.Fprintln(w, versions[0].Version)
		return err
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(versions)
}

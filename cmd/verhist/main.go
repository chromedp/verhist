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
	platform := flag.String("platform", "win64", "platform")
	channel := flag.String("channel", "stable", "channel")
	flag.Parse()
	if err := run(context.Background(), os.Stdout, *latest, *platform, *channel); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, w io.Writer, latest bool, platform, channel string) error {
	versions, err := verhist.Versions(ctx, platform, channel)
	switch {
	case err != nil:
		return err
	case len(versions) < 1:
		return verhist.ErrNoVersionsReturned
	}
	if latest {
		_, err := fmt.Fprintf(w, "%s\n", versions[0].Version)
		return err
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(versions)
}

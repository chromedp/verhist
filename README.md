# omahaproxy

Package `omahaproxy` provides a simple client to retrieve data from [Omaha
Proxy][omahaproxy].

[omahaproxy]: https://omahaproxy.appspot.com

## Example

```go
// _example/example.go
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/chromedp/omahaproxy"
)

func main() {
	platform := "linux"
	switch runtime.GOOS {
	case "windows":
		platform = "win"
		if runtime.GOARCH == "amd64" {
			platform += "64"
		}
	case "darwin":
		platform = "mac"
		if runtime.GOARCH == "aarch64" {
			platform += "_arm64"
		}
	}
	verbose := flag.Bool("v", false, "verbose")
	osstr := flag.String("os", platform, "os")
	channel := flag.String("channel", "stable", "channel")
	flag.Parse()
	if err := run(context.Background(), *verbose, *osstr, *channel); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, verbose bool, os, channel string) error {
	// enable verbose
	var opts []omahaproxy.Option
	if verbose {
		opts = append(opts, omahaproxy.WithLogf(fmt.Printf))
	}
	// create client
	cl := omahaproxy.New(opts...)
	// retrieve recent
	releases, err := cl.Recent(ctx)
	if err != nil {
		return err
	}
	for _, release := range releases {
		fmt.Printf("os: %s channel: %s version: %s\n", release.OS, release.Channel, release.Version)
	}
	// show latest
	ver, err := cl.Latest(ctx, os, channel)
	if err != nil {
		return err
	}
	fmt.Printf("latest: %s\n", ver)
	return nil
}
```

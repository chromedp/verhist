# verhist

Package `verhist` provides a simple client to retrieve the latest release
versions of Chrome using the [version history API][verhist].

[verhist]: https://developer.chrome.com/docs/web-platform/versionhistory/guide

Can also be used to build the latest user agent for Chrome.

## Example

```go
package verhist_test

import (
	"context"
	"fmt"

	"github.com/chromedp/verhist"
)

func Example() {
	userAgent, err := verhist.UserAgent(context.Background(), "linux", "stable")
	if err != nil {
		panic(err)
	}
	fmt.Println(userAgent)
}
```

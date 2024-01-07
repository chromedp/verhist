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
	fmt.Println(userAgent != "")
	// Output:
	// true
}

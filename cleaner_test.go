package cleaner_test

import (
	"fmt"

	"github.com/jabolopes/go-cleaner"
)

func ExampleGroup() {
	g, cleanup := cleaner.New()
	defer cleanup()

	g.Add(func() { fmt.Printf("hello") })

	// Output: hello
}

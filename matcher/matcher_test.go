package matcher_test

import (
	"fmt"

	"github.com/alextanhongpin/pkg/matcher"
)

func Example() {
	header := matcher.Match("Bearer")
	switch {
	case header.Is("Bearer"):
		fmt.Println("is bearer")
	case header.Is("Basic"):
		fmt.Println("is basic")
	default:
		fmt.Println("none")
	}
}

package matcher

import (
	"strings"
)

/*
const (
	Bearer = "bearer"
	Basic  = "basic"
)

func main() {
	header := Match("Bearer")
	switch {
	case header.Is("Bearer"):
		fmt.Println("is bearer")
	case header.Is("Basic"):
		fmt.Println("is basic")
	default:
		fmt.Println("none")
	}
}
*/

type Match string

// Is checks if the strings are similar. Capitalization does not matter.
func (m Match) Is(s string) bool {
	return strings.EqualFold(string(m), s)
}

// IsStrict checks if the strings are equal.
func (m Match) IsStrict(s string) bool {
	return string(m) == s
}

// IsEmpty checks if the string is not an empty space.
func (m Match) IsEmpty() bool {
	return len(strings.TrimSpace(string(m))) == 0
}

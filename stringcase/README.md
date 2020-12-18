## Another implementation

```go
package main

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	wordSeparatorRe       *regexp.Regexp
	underscoreRe          *regexp.Regexp
	duplicateUnderscoreRe *regexp.Regexp
	duplicateSpaceRe      *regexp.Regexp
	nonWordsRe            *regexp.Regexp
)

func init() {
	// A word can be separated if the current character is lowercase and the
	// next character that follows it is uppercase, a number or a space.
	wordSeparatorRe = regexp.MustCompile(`([a-z])([A-Z0-9 ])`)
	underscoreRe = regexp.MustCompile(`_\w`)
	nonWordsRe = regexp.MustCompile(`\W+`)
	duplicateUnderscoreRe = regexp.MustCompile(`_+`)
	duplicateSpaceRe = regexp.MustCompile(` +`)
}

func SnakeCase(in string) string {
	in = nonWordsRe.ReplaceAllString(in, " ")
	in = wordSeparatorRe.ReplaceAllString(in, "${1} ${2}")
	in = duplicateSpaceRe.ReplaceAllString(in, " ")
	in = strings.TrimSpace(in)
	in = strings.Replace(in, " ", "_", -1)
	return strings.ToLower(in)
}

func CamelCase(in string) string {
	in = SnakeCase(in)

	out := underscoreRe.ReplaceAllStringFunc(in, func(str string) string {
		return strings.ToUpper(str[1:])
	})
	return out
}

func KebabCase(in string) string {
	in = SnakeCase(in)
	return strings.Replace(in, "_", "-", -1)
}

func PascalCase(in string) string {
	in = CamelCase(in)
	return strings.ToUpper(in[:1]) + in[1:]
}

func main() {
	Example("SnakeCase", SnakeCase)
	Example("CamelCase", CamelCase)
	Example("KebabCase", KebabCase)
	Example("PascalCase", PascalCase)
}

func Example(name string, algo func(string) string) {
	tests := []string{
		"user_service",
		"party pooper",
		"THE AMAZING SPIDER-MAN",
		"A!@#$SAFS ridiculous",
		"property changer",
		"userID",
		"USER ID.",
		"created+at",
		"!created+at!",
		"this--is--a--slug",
		"address.home",
		"ZamZam Alakazam",
		"HelloWorld",
		"user 1000",
		"user1000",
		"helloWorld",
		"this is.   amazing	",
	}
	fmt.Println(name)
	for _, t := range tests {
		fmt.Println(algo(t))
	}
	fmt.Println()
}
```

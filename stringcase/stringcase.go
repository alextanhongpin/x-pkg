package stringcase

import (
	"log"
	"regexp"
	"strings"
)

var (
	splitRe *regexp.Regexp
	camelRe *regexp.Regexp
	trailRe *regexp.Regexp
)

func init() {
	var err error
	splitRe, err = regexp.Compile("(i?)[^a-z0-9]+[a-z0-9]")
	if err != nil {
		log.Fatal(err)
	}
	camelRe, err = regexp.Compile("[a-z][A-Z]")
	if err != nil {
		log.Fatal(err)
	}
	trailRe, err = regexp.Compile("(?i)^[^a-z0-9]|[^a-z0-9]$")
	if err != nil {
		log.Fatal(err)
	}
}

func CamelCase(s string) string {
	s = camelRe.ReplaceAllStringFunc(s, func(s string) string {
		return s[:1] + "-" + strings.ToLower(s[1:])
	})
	s = trailRe.ReplaceAllString(s, "")
	s = strings.ToLower(s)
	s = splitRe.ReplaceAllStringFunc(s, func(s string) string {
		return strings.ToUpper(s[len(s)-1:])
	})
	return s
}

func SnakeCase(s string) string {
	if s == "" {
		return s
	}
	camel := CamelCase(s)
	snake := camelRe.ReplaceAllStringFunc(camel, func(s string) string {
		return s[:1] + "_" + strings.ToLower(s[1:])
	})
	return strings.ToLower(snake)
}

func KebabCase(s string) string {

	if s == "" {
		return s
	}
	camel := CamelCase(s)
	kebab := camelRe.ReplaceAllStringFunc(camel, func(s string) string {
		return s[:1] + "-" + strings.ToLower(s[1:])
	})
	return strings.ToLower(kebab)
}

func PascalCase(s string) string {
	if s == "" {
		return s
	}
	camel := CamelCase(s)
	return strings.ToUpper(camel[:1]) + camel[1:]
}

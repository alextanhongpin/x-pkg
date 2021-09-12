package stringcase

import (
	"regexp"
	"strings"
)

var (
	sepRe   *regexp.Regexp
	camelRe *regexp.Regexp
)

func init() {
	sepRe = regexp.MustCompile(`(?i)([^a-z0-9]+)`)
	camelRe = regexp.MustCompile(`([a-z])([A-Z0-9])`)
}

type StringCaser func(string) string

var KebabCase = Pipe(identity, toKebab)
var PascalCase = Pipe(identity, toKebab, kebabToPascal)
var CamelCase = Pipe(identity, toKebab, kebabToPascal, lowerFirst)
var SnakeCase = Pipe(identity, toKebab, kebabToSnake)

func kebabToPascal(next StringCaser) StringCaser {
	return func(s string) string {
		parts := strings.Split(s, "-")
		result := make([]string, len(parts))
		for i, part := range parts {
			if i == len(parts)-1 {
				result[i] = UpperCommonInitialism(part)
			} else {
				result[i] = UpperFirst(part)
			}
		}
		s = strings.Join(result, "")
		return next(s)
	}
}

func kebabToSnake(next StringCaser) StringCaser {
	return func(s string) string {
		s = strings.ReplaceAll(s, "-", "_")
		return next(s)
	}
}

func toKebab(next StringCaser) StringCaser {
	return func(s string) string {
		s = camelRe.ReplaceAllString(s, "$1 $2")
		s = sepRe.ReplaceAllString(s, " ")
		s = strings.TrimSpace(s)
		s = strings.ReplaceAll(s, " ", "-")
		s = strings.ToLower(s)
		return next(s)
	}
}

func lowerFirst(next StringCaser) StringCaser {
	return func(s string) string {
		s = LowerFirst(s)
		return next(s)
	}
}

func identity(s string) string {
	return s
}

func Pipe(head StringCaser, next ...func(StringCaser) StringCaser) func(string) string {
	for i := range next {
		head = next[len(next)-i-1](head)
	}
	return head
}

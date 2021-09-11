package stringcase

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	splitRe  *regexp.Regexp
	camelRe  *regexp.Regexp
	trailRe  *regexp.Regexp
	repeatRe *regexp.Regexp
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
	trailRe, err = regexp.Compile("(?i)^[^a-z0-9]+|[^a-z0-9]+$")
	if err != nil {
		log.Fatal(err)
	}
	repeatRe, err = regexp.Compile("(?i)[^a-z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
}

type StringCaser func(string) string

var KebabCase = Pipe(identity, camelToKebab, replaceRepeatingSpecialCharactersWithHyphen, lowerAll)
var PascalCase = Pipe(identity, normalizeToKebab, kebabToPascal)
var CamelCase = Pipe(identity, normalizeToKebab, kebabToPascal, lowerFirst)
var SnakeCase = Pipe(identity, normalizeToKebab, kebabToSnake)

func normalizeToKebab(next StringCaser) StringCaser {
	return func(s string) string {
		s = KebabCase(s)
		return next(s)
	}
}

func camelToKebab(next StringCaser) StringCaser {
	return func(s string) string {
		s = camelRe.ReplaceAllStringFunc(s, func(s string) string {
			h, size := utf8.DecodeRuneInString(s)
			if size == 0 {
				return ""
			}
			t, size := utf8.DecodeRuneInString(s[size:])
			if size == 0 {
				return fmt.Sprintf("%c", h)
			}
			return fmt.Sprintf("%c-%c", h, t)
		})
		return next(s)
	}
}

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

func replaceRepeatingSpecialCharactersWithHyphen(next StringCaser) StringCaser {
	return func(s string) string {
		s = repeatRe.ReplaceAllString(s, "-")
		s = trailRe.ReplaceAllString(s, "")
		return next(s)
	}
}

func lowerAll(next StringCaser) StringCaser {
	return func(s string) string {
		return next(strings.ToLower(s))
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

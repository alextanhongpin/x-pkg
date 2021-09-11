package stringcase

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// UpperFirst uppers the first character in a string.
func UpperFirst(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if size == 0 {
		return ""
	}
	return fmt.Sprintf("%c%s", unicode.ToUpper(r), s[size:])
}

// LowerFirst lowers the first character in a string.
func LowerFirst(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if size == 0 {
		return ""
	}
	return fmt.Sprintf("%c%s", unicode.ToLower(r), s[size:])
}

// UpperCommonInitialism uppers the character in a string based on common initialism.
func UpperCommonInitialism(s string) string {
	isCommonInitialisms := commonInitialisms[strings.ToUpper(s)]
	if isCommonInitialisms {
		return strings.ToUpper(s)
	}
	return UpperFirst(s)
}

// https://github.com/golang/lint/blob/83fdc39ff7b56453e3793356bcff3070b9b96445/lint.go#L770-L809
// commonInitialisms is a set of common initialisms.
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
var commonInitialisms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}

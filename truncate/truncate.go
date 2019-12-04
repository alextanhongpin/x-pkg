// Package truncate provides methods to truncate sentence and paragraph to a
// more optimal length.
package truncate

import (
	"strings"
)

// Sentence breaks the sentence into half.
func Sentence(text string) string {
	tokens := strings.Fields(text)
	mid := len(tokens) / 2
	sentence := strings.Join(tokens[:mid], " ")
	return strings.Join([]string{sentence, "..."}, " ")
}

// Paragraph breaks the sentence into the desired length.
func Paragraph(text string, n int) string {
	text = strings.TrimSpace(text)
	if len(text) < n {
		return text
	}
	sentences := strings.Split(text, ".")
	var chars, i int
	for chars < n {
		chars += len(strings.TrimSpace(sentences[i]))
		i++
	}
	result := append(sentences[:i-1], Sentence(sentences[i-1]))
	return strings.Join(result, ". ")
}

package profane_test

import (
	"testing"

	"github.com/alextanhongpin/pkg/profane"
)

func init() {
	profane.Add("six")
	profane.Add("hello")
}

func TestReplaceGarbled(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"", ""},
		{" ", " "},
		{"hello world", "$@!#% world"},
		{"HELLO world", "$@!#% world"},
	}
	for _, tt := range tests {
		got := profane.ReplaceGarbled(tt.in)
		if got != tt.out {
			t.Fatalf("expected %s, got %s", tt.out, got)
		}
	}
}

func TestReplaceStars(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"", ""},
		{" ", " "},
		{"six", "s*x"},
		{"SIX", "S*X"},
		{"hello world", "h***o world"},
		{"HELLO WORLD", "H***O WORLD"},
		{"HELLO HELLO", "H***O H***O"},
	}

	for _, tt := range tests {
		got := profane.ReplaceStars(tt.in)
		if got != tt.out {
			t.Fatalf("expected %s, got %s", tt.out, got)
		}
	}
}

func TestReplaceVowels(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"", ""},
		{" ", " "},
		{"hello world", "h*ll* world"},
		{"HELLO WORLD", "H*LL* WORLD"},
		{"HELLO HELLO", "H*LL* H*LL*"},
	}
	for _, tt := range tests {
		got := profane.ReplaceVowels(tt.in)
		if got != tt.out {
			t.Fatalf("expected %s, got %s", tt.out, got)
		}
	}
}

func TestReplaceNonConsonants(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"", ""},
		{" ", " "},
		{"hello world", "h*ll* world"},
		{"HELLO WORLD", "H*LL* WORLD"},
		{"HELLO HELLO", "H*LL* H*LL*"},
	}
	for _, tt := range tests {
		got := profane.ReplaceNonConsonants(tt.in)
		if got != tt.out {
			t.Fatalf("expected %s, got %s", tt.out, got)
		}
	}
}

func TestReplaceCustom(t *testing.T) {
	tests := []struct {
		in     string
		out    string
		custom string
	}{
		{"", "", ""},
		{" ", " ", ""},
		{"hello world", "[CENSORED] world", "[CENSORED]"},
		{"hello HELLO", "[CENSORED] [CENSORED]", "[CENSORED]"},
	}
	for _, tt := range tests {
		got := profane.ReplaceCustom(tt.in, tt.custom)
		if got != tt.out {
			t.Fatalf("expected %s, got %s", tt.out, got)
		}
	}
}

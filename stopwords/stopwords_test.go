package stopwords_test

import (
	"fmt"
	"testing"

	"github.com/alextanhongpin/pkg/stopwords"
)

func Example() {
	fmt.Println(stopwords.Has("i"))
	// Output: true

	stopwords.Add("xyz", "zyx")
	fmt.Println(stopwords.Has("xyz"))
	// Output: true
}

func Example_new() {
	sw := stopwords.New()
	sw.Add("xyz", "yzx")
	fmt.Println(sw.Has("xyz"))
	// Output: true

	stopwords.Has("xyz")
	// Output: false
}

func TestHas(t *testing.T) {
	expected := true
	actual := stopwords.Has("i")
	if expected != actual {
		t.Fatalf("expected %t, got %t", expected, actual)
	}
}

func TestNew(t *testing.T) {
	sw := stopwords.New()
	sw.Add("xyz")

	expected := true
	actual := sw.Has("xyz")
	if expected != actual {
		t.Fatalf("expected %t, got %t", expected, actual)
	}

	expected = false
	actual = stopwords.Has("xyz")
	if expected != actual {
		t.Fatalf("expected %t, got %t", expected, actual)
	}
}

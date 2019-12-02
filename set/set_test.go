package set_test

import (
	"testing"

	"github.com/alextanhongpin/pkg/set"
)

func TestSet(t *testing.T) {
	s := set.New("a", "b", "c")
	s.Add("d", "e")

	if size := s.Size(); size != 5 {
		t.Fatalf("expected %d, got %d", 5, size)
	}

	s.Remove("a", "b")
	if exists := s.Has("a"); exists {
		t.Fatalf("expected %t, got %t", false, exists)
	}
}

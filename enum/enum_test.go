package enum_test

import (
	"testing"

	"github.com/alextanhongpin/pkg/enum"
)

func TestEnum(t *testing.T) {
	const (
		North = "north"
		West  = "west"
		South = "south"
		East  = "east"
	)

	directions := enum.New(West, South, East)
	directions.Add(North, North)
	if got := directions.Has(North); !got {
		t.Fatalf("expected %t, got %t", true, got)
	}

	if got := directions.Has("north"); !got {
		t.Fatalf("expected %t, got %t", true, got)
	}

	if got := directions.Has("North"); got {
		t.Fatalf("expected %t, got %t", false, got)
	}
}

func TestStringEnum(t *testing.T) {
	const (
		North = "north"
		West  = "west"
		South = "south"
		East  = "east"
	)

	directions := enum.NewString(West, South, East)
	directions.Add(North, North)
	if got := directions.Has(North); !got {
		t.Fatalf("expected %t, got %t", true, got)
	}

	if got := directions.Has("north"); !got {
		t.Fatalf("expected %t, got %t", true, got)
	}

	if got := directions.Has("North"); !got {
		t.Fatalf("expected %t, got %t", true, got)
	}

	if got := directions.HasStrict("North"); got {
		t.Fatalf("expected %t, got %t", false, got)
	}
}

package enum_test

import (
	"fmt"
	"testing"

	"github.com/alextanhongpin/pkg/enum"
)

func Example() {
	const (
		North = "north"
		West  = "west"
		South = "south"
		East  = "east"
	)
	directions := enum.New(North, West, South, East)
	fmt.Println(
		directions.Has(North),
		directions.Has("north"),
		directions.Has("North"),
	)
}

func Example_int() {
	directions := enum.New(0, 1, 2, 3)
	fmt.Println(
		directions.Has(1),
		directions.Has("2"),
		directions.Has(2),
	)
}

func Example_role() {
	const (
		User           = "user"
		Admin          = "admin"
		RestrictedUser = "restricted_user"
	)
	roles := enum.New(User, Admin)
	fmt.Println(roles.Has("Admin"))
	fmt.Println(roles.Has(RestrictedUser))
	roles.Add(RestrictedUser)
	fmt.Println(roles.Has(RestrictedUser))
}

func Example_new_string() {
	const (
		User  = "user"
		Admin = "admin"
	)
	roles := enum.NewString(User, Admin)
	fmt.Println(roles.Has("Admin"))
	fmt.Println(roles.Has("uSeR"))
}

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

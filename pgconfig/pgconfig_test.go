package pgconfig_test

import (
	"fmt"
	"testing"

	"github.com/alextanhongpin/pkg/pgconfig"
)

func Example() {
	cfg := pgconfig.New()
	fmt.Println(cfg.String())
	// Output: dbname='postgres' host='localhost' password='postgres' port='5432' sslmode='disable' user='root'
}

func TestConfig(t *testing.T) {
	got := pgconfig.New().String()
	expected := `dbname='postgres' host='localhost' password='postgres' port='5432' sslmode='disable' user='root'`

	if expected != got {
		t.Fatalf("expected %s, got %s", expected, got)
	}
}

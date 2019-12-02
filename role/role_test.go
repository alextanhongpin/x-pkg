package role_test

import (
	"fmt"
	"testing"

	"github.com/alextanhongpin/pkg/role"
	"github.com/alextanhongpin/pkg/set"
)

func Example() {
	var roles = role.Roles{
		"admin": set.New("read:books", "delete:books", "create:books"),
		"user":  set.New("read:books"),
	}

	fmt.Println("read:books", roles.Can("read:books"))
	fmt.Println("create:books", roles.Can("create:books"))
}

func BenchmarkRoleSet(b *testing.B) {
	var roles = Roles{
		"admin": set.New("read:books", "delete:books", "create:books"),
		"user":  set.New("read:books"),
	}
	for n := 0; n < b.N; n++ {
		roles.Can("read:books")
	}
}
func BenchmarkRoleSlice(b *testing.B) {
	var roles = role.Roles{
		"admin": []string{"read:books", "delete:books", "create:books"},
		"user":  []string{"read:books"},
	}
	for n := 0; n < b.N; n++ {
		roles.Can("read:books")
	}
}

type Roles map[string]set.Set

func (r Roles) Can(target string) []string {
	var result []string
	for role, scopes := range r {
		if scopes.Has(target) {
			result = append(result, role)
		}
	}
	return result
}

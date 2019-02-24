package enum

import (
	"strings"
)

type (
	Enum interface {
		Has(interface{}) bool
		Add(interface{})
	}

	EnumImpl struct {
		value map[interface{}]struct{}
	}
)

func New(val interface{}, values ...interface{}) *EnumImpl {
	value := make(map[interface{}]struct{})
	value[val] = struct{}{}
	for _, val := range values {
		value[val] = struct{}{}
	}

	return &EnumImpl{value}
}

func (e *EnumImpl) Has(enum interface{}) bool {
	_, exist := e.value[enum]
	return exist
}

func (e *EnumImpl) Add(enum interface{}) {
	e.value[enum] = struct{}{}
}

type (
	StringEnum interface {
		Has(string) bool
		HasStrict(string) bool
		Add(string)
	}

	StringEnumImpl struct {
		value []string
	}
)

func NewString(val string, values ...string) *StringEnumImpl {
	value := make([]string, len(values)+1)
	value[0] = val
	for i, val := range values {
		value[i+1] = val
	}
	return &StringEnumImpl{value}
}

func (s *StringEnumImpl) Has(enum string) bool {
	for _, val := range s.value {
		if eq := strings.EqualFold(val, enum); eq {
			return eq
		}

	}
	return false
}

func (s *StringEnumImpl) HasStrict(enum string) bool {
	for _, val := range s.value {
		if eq := val == enum; eq {
			return eq
		}
	}
	return false
}

/*

func main() {
	const (
		North = "north"
		West  = "west"
		South = "south"
		East  = "east"
	)

	{
		directions := New(North, West, South, East)
		fmt.Println(
			directions.Has(North),
			directions.Has("north"),
			directions.Has("North"),
		)
	}
	{
		directions := New(0, 1, 2, 3)
		fmt.Println(
			directions.Has(1),
			directions.Has("2"),
			directions.Has(2),
		)
	}

	{
		const (
			User           = "user"
			Admin          = "admin"
			RestrictedUser = "restricted_user"
		)
		roles := New(User, Admin)
		fmt.Println(roles.Has("Admin"))
		fmt.Println(roles.Has(RestrictedUser))
		roles.Add(RestrictedUser)
		fmt.Println(roles.Has(RestrictedUser))
	}

	{
		const (
			User  = "user"
			Admin = "admin"
		)
		roles := NewString(User, Admin)
		fmt.Println(roles.Has("Admin"))
		fmt.Println(roles.Has("uSeR"))
	}
}

*/

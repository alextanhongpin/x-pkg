// Package enum allows creation of enum and comparing between enums.
package enum

import (
	"strings"
)

type (
	// Comparable represents the enum operations.
	Comparable interface {
		Has(interface{}) bool
		Add(interface{})
	}

	// Enum represents an enum type.
	Enum struct {
		value map[interface{}]struct{}
	}
)

// New returns a new Enum.
func New(val interface{}, values ...interface{}) *Enum {
	value := make(map[interface{}]struct{})
	value[val] = struct{}{}
	for _, val := range values {
		value[val] = struct{}{}
	}

	return &Enum{value}
}

// Has checks if the enum exists.
func (e *Enum) Has(enum interface{}) bool {
	_, exist := e.value[enum]
	return exist
}

// Add adds the enum to the existing enums.
func (e *Enum) Add(enums ...interface{}) {
	for _, enum := range enums {
		e.value[enum] = struct{}{}
	}
}

type (
	// StringsComparable represents an string enums.
	StringsComparable interface {
		Has(string) bool
		HasStrict(string) bool
		Add(string)
	}

	// StringEnum holds a list of string enums.
	StringEnum struct {
		value map[string]struct{}
	}
)

// NewString adds a string enum to the existing collections.
func NewString(val string, values ...string) *StringEnum {

	value := make(map[string]struct{})
	value[val] = struct{}{}
	for _, val := range values {
		value[val] = struct{}{}
	}
	return &StringEnum{value}
}

// Add adds the enum to the existing enums.
func (s *StringEnum) Add(enums ...string) {
	for _, enum := range enums {
		s.value[enum] = struct{}{}
	}
}

// Has checks if a string exists in the current enums, ignoring cases.
func (s *StringEnum) Has(enum string) bool {
	for val := range s.value {
		if eq := strings.EqualFold(val, enum); eq {
			return eq
		}
	}
	return false
}

// HasStrict checks if a string with matching cases exists.
func (s *StringEnum) HasStrict(enum string) bool {
	for val := range s.value {
		if eq := val == enum; eq {
			return eq
		}
	}
	return false
}

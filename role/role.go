// Package role makes it easy to compare role and scopes.
package role

// Roles holds a map of roles and its respective scopes.
type Roles map[string][]string

// Can returns the roles that has the given scope.
func (r Roles) Can(target string) (result []string) {
	for role, scopes := range r {
		for _, scope := range scopes {
			if scope == target {
				result = append(result, role)
				break
			}
		}
	}
	return
}

package role

type Roles map[string][]string

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

// func main() {
//         var roles = Roles{
//                 "admin": NewSet("read:books", "delete:books", "create:books"),
//                 "user":  NewSet("read:books"),
//         }
//
//         fmt.Println("read:books", roles.Can("read:books"))
//         fmt.Println("create:books", roles.Can("create:books"))
// }

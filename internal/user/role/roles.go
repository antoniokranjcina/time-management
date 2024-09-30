package role

// Role represents the string type for user roles.
type Role string

// Constants representing different user roles as strings.
const (
	SuperAdmin Role = "super_admin"
	Admin      Role = "admin"
	Manager    Role = "manager"
	Employee   Role = "employee"
)

// validRoles is a map of valid roles for lookup.
var validRoles = map[Role]bool{
	SuperAdmin: true,
	Admin:      true,
	Manager:    true,
	Employee:   true,
}

// String method provides the string representation of a Role (optional, since Role is already a string).
func (r Role) String() string {
	return string(r)
}

func Strings(roles ...Role) []string {
	roleStrings := make([]string, len(roles))
	for i, r := range roles {
		roleStrings[i] = r.String()
	}
	return roleStrings
}

// IsValid checks if the role is a valid Role value.
func (r Role) IsValid() bool {
	return validRoles[r]
}

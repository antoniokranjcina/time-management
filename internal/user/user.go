package user

const TableName = "users"

type User struct {
	Id           string `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	PasswordHash string `json:"password_hash"`
	CreatedAt    uint64 `json:"created_at"`
	Active       bool   `json:"active"`
}

// NewAdmin Factory method to create an Admin
func NewAdmin(id, firstName, lastName, email, password string, createdAt uint64, active bool) *User {
	return &User{
		Id:           id,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		Role:         "admin",
		PasswordHash: password,
		CreatedAt:    createdAt,
		Active:       active,
	}
}

// NewModerator Factory method to create a Moderator
func NewModerator(id, firstName, lastName, email, password string, createdAt uint64, active bool) *User {
	return &User{
		Id:           id,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		Role:         "moderator",
		PasswordHash: password,
		CreatedAt:    createdAt,
		Active:       active,
	}
}

// NewEmployee Factory method to create an Employee
func NewEmployee(id, firstName, lastName, email, password string, createdAt uint64, active bool) *User {
	return &User{
		Id:           id,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		Role:         "employee",
		PasswordHash: password,
		CreatedAt:    createdAt,
		Active:       active,
	}
}

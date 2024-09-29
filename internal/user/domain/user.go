package domain

import "time-management/internal/user/role"

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
		Role:         role.Admin.String(),
		PasswordHash: password,
		CreatedAt:    createdAt,
		Active:       active,
	}
}

// NewManager Factory method to create a Moderator
func NewManager(id, firstName, lastName, email, password string, createdAt uint64, active bool) *User {
	return &User{
		Id:           id,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		Role:         role.Manager.String(),
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
		Role:         role.Employee.String(),
		PasswordHash: password,
		CreatedAt:    createdAt,
		Active:       active,
	}
}

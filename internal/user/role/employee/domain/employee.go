package domain

import (
	"time-management/internal/user/domain"
)

type Employee struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	CreatedAt uint64 `json:"created_at"`
	Active    bool   `json:"active"`
}

func MapUserToEmployee(user *domain.User) *Employee {
	return &Employee{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		Active:    user.Active,
	}
}

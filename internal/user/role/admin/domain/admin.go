package domain

import "time-management/internal/user/domain"

type Admin struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	CreatedAt uint64 `json:"created_at"`
	Active    bool   `json:"active"`
}

func MapUserToAdmin(user *domain.User) *Admin {
	return &Admin{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		Active:    user.Active,
	}
}

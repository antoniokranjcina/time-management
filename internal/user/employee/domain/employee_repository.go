package domain

import "time-management/internal/user"

type EmployeeRepository interface {
	Save(employee *user.User) (*user.User, error)
	GetAll() ([]user.User, error)
	GetById(id string) (*user.User, error)
	Update(id, firstName, lastName string) (*user.User, error)
	ChangePassword(id, password string) error
	ChangeEmail(id, email string) error
	ToggleStatus(id string, status bool) (bool, error)
	Delete(id string) error
}

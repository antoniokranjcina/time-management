package domain

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetAllWithRole(ctx context.Context, role string) ([]User, error)
	GetByIdWithRole(ctx context.Context, id, role string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, id, firstName, lastName string) (*User, error)
	ChangePassword(ctx context.Context, id, password string) error
	ChangeEmail(ctx context.Context, id, email string) error
	ToggleStatus(ctx context.Context, id string, status bool) (bool, error)
	Delete(ctx context.Context, id string) error
}

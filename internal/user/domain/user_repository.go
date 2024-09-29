package domain

type UserRepository interface {
	Save(employee *User) (*User, error)
	GetAllWithRole(role string) ([]User, error)
	GetByIdWithRole(id, role string) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(id, firstName, lastName string) (*User, error)
	ChangePassword(id, password string) error
	ChangeEmail(id, email string) error
	ToggleStatus(id string, status bool) (bool, error)
	Delete(id string) error
}

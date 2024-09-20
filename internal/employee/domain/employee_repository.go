package domain

type EmployeeRepository interface {
	Save(employee *Employee) (*Employee, error)
	GetAll() ([]Employee, error)
	GetById(id string) (*Employee, error)
	Update(id, firstName, lastName string) (*Employee, error)
	ChangePassword(id, password string) error
	ChangeEmail(id, email string) error
	ToggleStatus(id string, status bool) (bool, error)
	Delete(id string) error
}

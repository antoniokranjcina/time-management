package domain

type Employee struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt uint64 `json:"created_at"`
	Active    bool   `json:"active"`
}

// NewEmployee Factory method to create a Employee
func NewEmployee(id, firstName, lastName, email, password string, createdAt uint64, active bool) *Employee {
	return &Employee{
		Id:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
		CreatedAt: createdAt,
		Active:    active,
	}
}

package command

import (
	"github.com/google/uuid"
	"time"
	"time-management/internal/employees/domain"
	"time-management/internal/employees/util"
)

type CreateEmployeeCommand struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type CreateEmployeeHandler struct {
	Repo domain.EmployeeRepository
}

func (h *CreateEmployeeHandler) Handle(cmd CreateEmployeeCommand) (*domain.Employee, error) {
	// Validate logic
	if cmd.FirstName == "" {
		return nil, domain.NewValidationError(domain.ErrFirstNameTooShort)
	}
	if cmd.LastName == "" {
		return nil, domain.NewValidationError(domain.ErrLastNameTooShort)
	}
	if !util.EmailRegex.MatchString(cmd.Email) {
		return nil, domain.NewValidationError(domain.ErrEmailWrongFormat)
	}
	if len(cmd.Password) < 6 {
		return nil, domain.NewValidationError(domain.ErrPasswordTooShort)
	}

	// Create the domain entity
	employee := domain.NewEmployee(
		uuid.New().String(),
		cmd.FirstName,
		cmd.LastName,
		cmd.Email,
		cmd.Password,
		uint64(time.Now().Unix()),
		true,
	)

	// Save it through repository
	createdEmployee, err := h.Repo.Save(employee)
	if err != nil {
		return nil, err
	}

	return createdEmployee, nil
}

package command

import (
	"github.com/google/uuid"
	"net/mail"
	"time"
	sharedUtil "time-management/internal/shared/util"
	"time-management/internal/user"
	"time-management/internal/user/employee/domain"
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
	if cmd.FirstName == "" {
		return nil, sharedUtil.NewValidationError(domain.ErrFirstNameTooShort)
	}
	if cmd.LastName == "" {
		return nil, sharedUtil.NewValidationError(domain.ErrLastNameTooShort)
	}
	if _, err := mail.ParseAddress(cmd.Email); err != nil {
		return nil, sharedUtil.NewValidationError(domain.ErrEmailWrongFormat)
	}
	if len(cmd.Password) < 6 {
		return nil, sharedUtil.NewValidationError(domain.ErrPasswordTooShort)
	}

	employee := user.NewEmployee(
		uuid.New().String(),
		cmd.FirstName,
		cmd.LastName,
		cmd.Email,
		cmd.Password,
		uint64(time.Now().Unix()),
		true,
	)

	createdUser, err := h.Repo.Save(employee)
	if err != nil {
		return nil, err
	}

	createdEmployee := domain.MapUserToEmployee(createdUser)

	return createdEmployee, nil
}

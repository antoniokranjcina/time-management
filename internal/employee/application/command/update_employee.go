package command

import (
	"time-management/internal/employee/domain"
	sharedUtil "time-management/internal/shared/util"
)

type UpdateEmployeeCommand struct {
	Id        string
	FirstName string
	LastName  string
}

type UpdateEmployeeHandler struct {
	Repo domain.EmployeeRepository
}

func (h *UpdateEmployeeHandler) Handle(cmd UpdateEmployeeCommand) (*domain.Employee, error) {
	// Validation
	if cmd.FirstName == "" {
		return nil, sharedUtil.NewValidationError(domain.ErrFirstNameTooShort)
	}
	if cmd.LastName == "" {
		return nil, sharedUtil.NewValidationError(domain.ErrLastNameTooShort)
	}

	// Update the domain entity through repository
	updatedEmployee, err := h.Repo.Update(cmd.Id, cmd.FirstName, cmd.LastName)
	if err != nil {
		return nil, err
	}

	return updatedEmployee, nil
}

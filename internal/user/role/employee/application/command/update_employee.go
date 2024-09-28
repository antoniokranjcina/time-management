package command

import (
	sharedUtil "time-management/internal/shared/util"
	"time-management/internal/user/domain"
	empDomain "time-management/internal/user/role/employee/domain"
)

type UpdateEmployeeCommand struct {
	Id        string
	FirstName string
	LastName  string
}

type UpdateEmployeeHandler struct {
	Repo domain.UserRepository
}

func (h *UpdateEmployeeHandler) Handle(cmd UpdateEmployeeCommand) (*empDomain.Employee, error) {
	if cmd.FirstName == "" {
		return nil, sharedUtil.NewValidationError(domain.ErrFirstNameTooShort)
	}
	if cmd.LastName == "" {
		return nil, sharedUtil.NewValidationError(domain.ErrLastNameTooShort)
	}

	updatedUser, err := h.Repo.Update(cmd.Id, cmd.FirstName, cmd.LastName)
	if err != nil {
		return nil, err
	}

	updatedEmployee := empDomain.MapUserToEmployee(updatedUser)

	return updatedEmployee, nil
}

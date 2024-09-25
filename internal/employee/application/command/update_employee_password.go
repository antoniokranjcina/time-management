package command

import (
	"time-management/internal/employee/domain"
	sharedUtil "time-management/internal/shared/util"
)

type UpdatePasswordCommand struct {
	Id       string
	Password string
}

type UpdatePasswordHandler struct {
	Repo domain.EmployeeRepository
}

func (h *UpdatePasswordHandler) Handle(cmd UpdatePasswordCommand) error {
	// Validation
	if len(cmd.Password) < 6 {
		return sharedUtil.NewValidationError(domain.ErrPasswordTooShort)
	}

	// Update the password through repository
	err := h.Repo.ChangePassword(cmd.Id, cmd.Password)
	if err != nil {
		return err
	}

	return nil
}

package command

import (
	sharedUtil "time-management/internal/shared/util"
	"time-management/internal/user/employee/domain"
)

type UpdatePasswordCommand struct {
	Id       string
	Password string
}

type UpdatePasswordHandler struct {
	Repo domain.EmployeeRepository
}

func (h *UpdatePasswordHandler) Handle(cmd UpdatePasswordCommand) error {
	if len(cmd.Password) < 6 {
		return sharedUtil.NewValidationError(domain.ErrPasswordTooShort)
	}

	err := h.Repo.ChangePassword(cmd.Id, cmd.Password)
	if err != nil {
		return err
	}

	return nil
}

package command

import (
	"time-management/internal/employee/domain"
	"time-management/internal/employee/util"
	sharedUtil "time-management/internal/shared/util"
)

type UpdateEmailCommand struct {
	Id    string
	Email string
}

type UpdateEmailHandler struct {
	Repo domain.EmployeeRepository
}

func (h *UpdateEmailHandler) Handle(cmd UpdateEmailCommand) error {
	// Validation
	if !util.EmailRegex.MatchString(cmd.Email) {
		return sharedUtil.NewValidationError(domain.ErrEmailWrongFormat)
	}

	// Update the email through repository
	err := h.Repo.ChangeEmail(cmd.Id, cmd.Email)
	if err != nil {
		return err
	}

	return nil
}

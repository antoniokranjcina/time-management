package command

import (
	"net/mail"
	sharedUtil "time-management/internal/shared/util"
	"time-management/internal/user/domain"
)

type UpdateEmailCommand struct {
	Id    string
	Email string
}

type UpdateEmailHandler struct {
	Repo domain.UserRepository
}

func (h *UpdateEmailHandler) Handle(cmd UpdateEmailCommand) error {
	if _, err := mail.ParseAddress(cmd.Email); err != nil {
		return sharedUtil.NewValidationError(domain.ErrEmailWrongFormat)
	}

	err := h.Repo.ChangeEmail(cmd.Id, cmd.Email)
	if err != nil {
		return err
	}

	return nil
}

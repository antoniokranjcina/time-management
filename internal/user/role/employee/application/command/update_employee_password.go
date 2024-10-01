package command

import (
	"golang.org/x/crypto/bcrypt"
	sharedUtil "time-management/internal/shared/util"
	"time-management/internal/user/domain"
)

type UpdatePasswordCommand struct {
	Id       string
	Password string
}

type UpdatePasswordHandler struct {
	Repo domain.UserRepository
}

func (h *UpdatePasswordHandler) Handle(cmd UpdatePasswordCommand) error {
	if len(cmd.Password) < 6 {
		return sharedUtil.NewValidationError(domain.ErrPasswordTooShort)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.MinCost)
	if err != nil {
		return domain.ErrInternalServer
	}

	err = h.Repo.ChangePassword(cmd.Id, string(hash))
	if err != nil {
		return err
	}

	return nil
}

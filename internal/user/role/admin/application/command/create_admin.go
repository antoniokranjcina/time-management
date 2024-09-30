package command

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"time"
	sharedUtil "time-management/internal/shared/util"
	"time-management/internal/user/domain"
	adminDomain "time-management/internal/user/role/admin/domain"
)

type CreateAdminCommand struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type CreateAdminHandler struct {
	Repo domain.UserRepository
}

func (h *CreateAdminHandler) Handle(cmd CreateAdminCommand) (*adminDomain.Admin, error) {
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

	hash, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.MinCost)
	if err != nil {
		return nil, domain.ErrInternalServer
	}

	admin := domain.NewAdmin(
		uuid.New().String(),
		cmd.FirstName,
		cmd.LastName,
		cmd.Email,
		string(hash),
		uint64(time.Now().Unix()),
		true,
	)

	createdUser, err := h.Repo.Save(admin)
	if err != nil {
		return nil, err
	}

	createdAdmin := adminDomain.MapUserToAdmin(createdUser)

	return createdAdmin, err
}

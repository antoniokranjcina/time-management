package command

import (
	"context"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"time"
	sharedUtil "time-management/internal/shared/util"
	"time-management/internal/user/domain"
	empDomain "time-management/internal/user/role/employee/domain"
)

type CreateEmployeeCommand struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type CreateEmployeeHandler struct {
	Repo domain.UserRepository
}

func (h *CreateEmployeeHandler) Handle(ctx context.Context, cmd CreateEmployeeCommand) (*empDomain.Employee, error) {
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

	employee := domain.NewEmployee(
		uuid.New().String(),
		cmd.FirstName,
		cmd.LastName,
		cmd.Email,
		string(hash),
		uint64(time.Now().Unix()),
		true,
	)

	createdUser, err := h.Repo.Create(ctx, employee)
	if err != nil {
		return nil, err
	}

	createdEmployee := empDomain.MapUserToEmployee(createdUser)

	return createdEmployee, nil
}

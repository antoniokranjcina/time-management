package command

import (
	"context"
	sharedUtil "time-management/internal/shared/util"
	"time-management/internal/user/domain"
	adminDomain "time-management/internal/user/role/admin/domain"
)

type UpdateAdminCommand struct {
	Id        string
	FirstName string
	LastName  string
}
type UpdateAdminHandler struct {
	Repo domain.UserRepository
}

func (h *UpdateAdminHandler) Handle(ctx context.Context, cmd UpdateAdminCommand) (*adminDomain.Admin, error) {
	if cmd.FirstName == "" {
		return nil, sharedUtil.NewValidationError(domain.ErrFirstNameTooShort)
	}
	if cmd.LastName == "" {
		return nil, sharedUtil.NewValidationError(domain.ErrLastNameTooShort)
	}

	updatedUser, err := h.Repo.Update(ctx, cmd.Id, cmd.FirstName, cmd.LastName)
	if err != nil {
		return nil, err
	}

	updatedAdmin := adminDomain.MapUserToAdmin(updatedUser)

	return updatedAdmin, nil
}

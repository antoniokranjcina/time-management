package command

import (
	"context"
	"time-management/internal/user/domain"
)

type DeleteAdminCommand struct {
	Id string
}

type DeleteAdminHandler struct {
	Repo domain.UserRepository
}

func (h *DeleteAdminHandler) Handle(ctx context.Context, cmd DeleteAdminCommand) error {
	err := h.Repo.Delete(ctx, cmd.Id)
	if err != nil {
		return err
	}

	return nil
}

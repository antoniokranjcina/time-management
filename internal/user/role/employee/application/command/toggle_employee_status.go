package command

import (
	"context"
	"time-management/internal/user/domain"
)

type ToggleStatusCommand struct {
	Id     string
	Active bool
}

type ToggleStatusHandler struct {
	Repo domain.UserRepository
}

func (h *ToggleStatusHandler) Handle(ctx context.Context, cmd ToggleStatusCommand) (bool, error) {
	newStatus, err := h.Repo.ToggleStatus(ctx, cmd.Id, cmd.Active)
	if err != nil {
		return cmd.Active, err
	}

	return newStatus, nil
}

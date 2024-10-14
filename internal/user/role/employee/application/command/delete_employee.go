package command

import (
	"context"
	"time-management/internal/user/domain"
)

type DeleteEmployeeCommand struct {
	Id string
}

type DeleteEmployeeHandler struct {
	Repo domain.UserRepository
}

func (h *DeleteEmployeeHandler) Handle(ctx context.Context, cmd DeleteEmployeeCommand) error {
	err := h.Repo.Delete(ctx, cmd.Id)
	if err != nil {
		return err
	}

	return nil
}

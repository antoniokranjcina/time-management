package command

import (
	"context"
	"time-management/internal/location/domain"
)

type DeleteLocationCommand struct {
	Id string
}

type DeleteLocationHandler struct {
	Repo domain.LocationRepository
}

func (h *DeleteLocationHandler) Handle(ctx context.Context, cmd DeleteLocationCommand) error {
	return h.Repo.Delete(ctx, cmd.Id)
}

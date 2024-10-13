package command

import (
	"context"
	"time-management/internal/location/domain"
	"time-management/internal/shared/util"
)

type UpdateLocationCommand struct {
	Id   string
	Name string
}

type UpdateLocationHandler struct {
	Repo domain.LocationRepository
}

func (h *UpdateLocationHandler) Handle(ctx context.Context, cmd UpdateLocationCommand) (*domain.Location, error) {
	// Validation logic
	if cmd.Name == "" || len(cmd.Name) >= 50 {
		return nil, util.NewValidationError(domain.ErrInvalidName)
	}

	// Update the domain entity through the repository
	updatedLocation, err := h.Repo.Update(ctx, cmd.Id, cmd.Name)
	if err != nil {
		return nil, err
	}

	return updatedLocation, nil
}

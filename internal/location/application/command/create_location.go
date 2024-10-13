package command

import (
	"context"
	"github.com/google/uuid"
	"time"
	"time-management/internal/location/domain"
	"time-management/internal/shared/util"
)

type CreateLocationCommand struct {
	Name string
}

type CreateLocationHandler struct {
	Repo domain.LocationRepository
}

func (h *CreateLocationHandler) Handle(ctx context.Context, cmd CreateLocationCommand) (*domain.Location, error) {
	if cmd.Name == "" || len(cmd.Name) >= 50 {
		return nil, util.NewValidationError(domain.ErrInvalidName)
	}

	location := domain.NewLocation(uuid.New().String(), cmd.Name, uint64(time.Now().Unix()))

	createdLocation, err := h.Repo.Create(ctx, location)
	if err != nil {
		return nil, err
	}

	return createdLocation, nil
}

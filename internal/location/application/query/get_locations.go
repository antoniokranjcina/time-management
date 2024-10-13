package query

import (
	"context"
	"time-management/internal/location/domain"
)

type GetLocationsHandler struct {
	Repo domain.LocationRepository
}

func (h *GetLocationsHandler) Handle(ctx context.Context) ([]domain.Location, error) {
	locations, err := h.Repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if locations == nil {
		return []domain.Location{}, nil
	}

	return locations, nil
}

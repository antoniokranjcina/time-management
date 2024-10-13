package query

import (
	"context"
	"time-management/internal/location/domain"
)

type GetLocationQuery struct {
	Id string
}

type GetLocationHandler struct {
	Repo domain.LocationRepository
}

func (h *GetLocationHandler) Handle(ctx context.Context, query GetLocationQuery) (*domain.Location, error) {
	location, err := h.Repo.GetById(ctx, query.Id)
	if err != nil {
		return nil, err
	}

	return location, nil
}

package queries

import "time-management/internal/locations/domain"

type GetLocationQuery struct {
	Id string
}

type GetLocationHandler struct {
	Repo domain.LocationRepository
}

func (h *GetLocationHandler) Handle(query GetLocationQuery) (*domain.Location, error) {
	location, err := h.Repo.GetById(query.Id)
	if err != nil {
		return nil, err
	}
	return location, nil
}

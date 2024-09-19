package query

import "time-management/internal/location/domain"

type GetLocationsHandler struct {
	Repo domain.LocationRepository
}

func (h *GetLocationsHandler) Handle() ([]domain.Location, error) {
	locations, err := h.Repo.GetAll()
	if err != nil {
		return nil, err
	}
	return locations, nil
}

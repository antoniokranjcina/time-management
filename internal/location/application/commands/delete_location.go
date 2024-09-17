package commands

import "time-management/internal/location/domain"

type DeleteLocationCommand struct {
	ID string
}

type DeleteLocationHandler struct {
	Repo domain.LocationRepository
}

func (h *DeleteLocationHandler) Handle(cmd DeleteLocationCommand) error {
	return h.Repo.Delete(cmd.ID)
}

package command

import "time-management/internal/location/domain"

type DeleteLocationCommand struct {
	Id string
}

type DeleteLocationHandler struct {
	Repo domain.LocationRepository
}

func (h *DeleteLocationHandler) Handle(cmd DeleteLocationCommand) error {
	return h.Repo.Delete(cmd.Id)
}

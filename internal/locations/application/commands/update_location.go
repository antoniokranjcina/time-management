package commands

import "time-management/internal/locations/domain"

type UpdateLocationCommand struct {
	Id   string
	Name string
}

type UpdateLocationHandler struct {
	Repo domain.LocationRepository
}

func (h *UpdateLocationHandler) Handle(cmd UpdateLocationCommand) (*domain.Location, error) {
	// Validation logic
	if cmd.Name == "" {
		return nil, domain.ErrInvalidName
	}

	// Update the domain entity through the repository
	updatedLocation, err := h.Repo.Update(cmd.Id, cmd.Name)
	if err != nil {
		return nil, err
	}

	return updatedLocation, nil
}

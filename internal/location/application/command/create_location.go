package command

import (
	"github.com/google/uuid"
	"time"
	"time-management/internal/location/domain"
)

type CreateLocationCommand struct {
	Name string
}

type CreateLocationHandler struct {
	Repo domain.LocationRepository
}

func (h *CreateLocationHandler) Handle(cmd CreateLocationCommand) (*domain.Location, error) {
	// Validation logic
	if cmd.Name == "" {
		return nil, domain.ErrInvalidName
	}

	// Create the domain entity
	location := domain.NewLocation(uuid.New().String(), cmd.Name, uint64(time.Now().Unix()))

	// Save it through the repository
	createdLocation, err := h.Repo.Save(location)
	if err != nil {
		return nil, err
	}

	return createdLocation, nil
}

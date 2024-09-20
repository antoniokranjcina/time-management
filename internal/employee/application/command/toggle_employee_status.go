package command

import "time-management/internal/employee/domain"

type ToggleStatusCommand struct {
	Id     string
	Active bool
}

type ToggleStatusHandler struct {
	Repo domain.EmployeeRepository
}

func (h *ToggleStatusHandler) Handle(cmd ToggleStatusCommand) (bool, error) {
	newStatus, err := h.Repo.ToggleStatus(cmd.Id, cmd.Active)
	if err != nil {
		return cmd.Active, err
	}

	return newStatus, nil
}

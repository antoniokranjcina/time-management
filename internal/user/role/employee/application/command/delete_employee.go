package command

import (
	"time-management/internal/user/domain"
)

type DeleteEmployeeCommand struct {
	Id string
}

type DeleteEmployeeHandler struct {
	Repo domain.UserRepository
}

func (h *DeleteEmployeeHandler) Handle(cmd DeleteEmployeeCommand) error {
	err := h.Repo.Delete(cmd.Id)
	if err != nil {
		return err
	}

	return nil
}

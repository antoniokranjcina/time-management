package command

import "time-management/internal/employees/domain"

type DeleteEmployeeCommand struct {
	Id string
}

type DeleteEmployeeHandler struct {
	Repo domain.EmployeeRepository
}

func (h *DeleteEmployeeHandler) Handle(cmd DeleteEmployeeCommand) error {
	return h.Repo.Delete(cmd.Id)
}

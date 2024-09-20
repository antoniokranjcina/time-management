package query

import "time-management/internal/employee/domain"

type GetEmployeesHandler struct {
	Repo domain.EmployeeRepository
}

func (h *GetEmployeesHandler) Handle() ([]domain.Employee, error) {
	employees, err := h.Repo.GetAll()
	if err != nil {
		return nil, err
	}
	return employees, nil
}

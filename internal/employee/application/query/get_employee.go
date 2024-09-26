package query

import "time-management/internal/employee/domain"

type GetEmployeeQuery struct {
	Id string
}

type GetEmployeeHandler struct {
	Repo domain.EmployeeRepository
}

func (h *GetEmployeeHandler) Handle(query GetEmployeeQuery) (*domain.Employee, error) {
	employee, err := h.Repo.GetById(query.Id)
	if err != nil {
		return nil, err
	}

	return employee, nil
}

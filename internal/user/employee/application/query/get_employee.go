package query

import (
	"time-management/internal/user/employee/domain"
)

type GetEmployeeQuery struct {
	Id string
}

type GetEmployeeHandler struct {
	Repo domain.EmployeeRepository
}

func (h *GetEmployeeHandler) Handle(query GetEmployeeQuery) (*domain.Employee, error) {
	user, err := h.Repo.GetById(query.Id)
	if err != nil {
		return nil, err
	}

	employee := domain.MapUserToEmployee(user)

	return employee, nil
}

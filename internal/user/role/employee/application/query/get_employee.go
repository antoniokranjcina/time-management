package query

import (
	"time-management/internal/user/domain"
	empDomain "time-management/internal/user/role/employee/domain"
)

type GetEmployeeQuery struct {
	Id string
}

type GetEmployeeHandler struct {
	Repo domain.UserRepository
}

func (h *GetEmployeeHandler) Handle(query GetEmployeeQuery) (*empDomain.Employee, error) {
	user, err := h.Repo.GetByIdWithRole(query.Id, "employee")
	if err != nil {
		return nil, err
	}

	employee := empDomain.MapUserToEmployee(user)

	return employee, nil
}

package query

import (
	"time-management/internal/user/domain"
	empDomain "time-management/internal/user/role/employee/domain"
)

type GetEmployeesHandler struct {
	Repo domain.UserRepository
}

func (h *GetEmployeesHandler) Handle() ([]empDomain.Employee, error) {
	users, err := h.Repo.GetAllWithRole("employee")
	if err != nil {
		return nil, err
	}

	var employees []empDomain.Employee
	for _, user := range users {
		employee := empDomain.MapUserToEmployee(&user)
		employees = append(employees, *employee)
	}

	return employees, nil
}

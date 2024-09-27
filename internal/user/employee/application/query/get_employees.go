package query

import (
	"time-management/internal/user/employee/domain"
)

type GetEmployeesHandler struct {
	Repo domain.EmployeeRepository
}

func (h *GetEmployeesHandler) Handle() ([]domain.Employee, error) {
	users, err := h.Repo.GetAll()
	if err != nil {
		return nil, err
	}

	var employees []domain.Employee
	for _, user := range users {
		employee := domain.MapUserToEmployee(&user)
		employees = append(employees, *employee)
	}

	return employees, nil
}

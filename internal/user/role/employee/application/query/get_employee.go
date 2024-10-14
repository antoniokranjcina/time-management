package query

import (
	"context"
	"time-management/internal/user/domain"
	"time-management/internal/user/role"
	empDomain "time-management/internal/user/role/employee/domain"
)

type GetEmployeeQuery struct {
	Id string
}

type GetEmployeeHandler struct {
	Repo domain.UserRepository
}

func (h *GetEmployeeHandler) Handle(ctx context.Context, query GetEmployeeQuery) (*empDomain.Employee, error) {
	user, err := h.Repo.GetByIdWithRole(ctx, query.Id, role.Employee.String())
	if err != nil {
		return nil, err
	}

	employee := empDomain.MapUserToEmployee(user)

	return employee, nil
}

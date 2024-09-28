package query

import (
	"time-management/internal/user/domain"
	adminDomain "time-management/internal/user/role/admin/domain"
)

type GetAdminQuery struct {
	Id string
}

type GetAdminHandler struct {
	Repo domain.UserRepository
}

func (h *GetAdminHandler) Handle(query GetAdminQuery) (*adminDomain.Admin, error) {
	user, err := h.Repo.GetByIdWithRole(query.Id, "admin")
	if err != nil {
		return nil, err
	}

	admin := adminDomain.MapUserToAdmin(user)

	return admin, nil
}

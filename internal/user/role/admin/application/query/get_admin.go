package query

import (
	"context"
	"time-management/internal/user/domain"
	"time-management/internal/user/role"
	adminDomain "time-management/internal/user/role/admin/domain"
)

type GetAdminQuery struct {
	Id string
}

type GetAdminHandler struct {
	Repo domain.UserRepository
}

func (h *GetAdminHandler) Handle(ctx context.Context, query GetAdminQuery) (*adminDomain.Admin, error) {
	user, err := h.Repo.GetByIdWithRole(ctx, query.Id, role.Admin.String())
	if err != nil {
		return nil, err
	}

	admin := adminDomain.MapUserToAdmin(user)

	return admin, nil
}

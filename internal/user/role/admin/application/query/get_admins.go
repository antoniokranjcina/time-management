package query

import (
	"context"
	"time-management/internal/user/domain"
	"time-management/internal/user/role"
	adminDomain "time-management/internal/user/role/admin/domain"
)

type GetAdminsHandler struct {
	Repo domain.UserRepository
}

func (h *GetAdminsHandler) Handle(ctx context.Context) ([]adminDomain.Admin, error) {
	users, err := h.Repo.GetAllWithRole(ctx, role.Admin.String())
	if err != nil {
		return nil, err
	}

	var admins []adminDomain.Admin
	for _, user := range users {
		admin := adminDomain.MapUserToAdmin(&user)
		admins = append(admins, *admin)
	}

	return admins, nil
}

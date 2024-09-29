package query

import (
	"time-management/internal/user/domain"
	"time-management/internal/user/role"
	adminDomain "time-management/internal/user/role/admin/domain"
)

type GetAdminsHandler struct {
	Repo domain.UserRepository
}

func (h *GetAdminsHandler) Handle() ([]adminDomain.Admin, error) {
	users, err := h.Repo.GetAllWithRole(role.Admin.String())
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

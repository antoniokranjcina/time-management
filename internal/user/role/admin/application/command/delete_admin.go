package command

import "time-management/internal/user/domain"

type DeleteAdminCommand struct {
	Id string
}

type DeleteAdminHandler struct {
	Repo domain.UserRepository
}

func (h *DeleteAdminHandler) Handle(cmd DeleteAdminCommand) error {
	err := h.Repo.Delete(cmd.Id)
	if err != nil {
		return err
	}

	return nil
}

package command

import (
	"time-management/internal/report/domain"
)

type DenyReportCommand struct {
	Id string
}

type DenyReportHandler struct {
	Repo domain.ReportRepository
}

func (h *DenyReportHandler) Handle(cmd DenyReportCommand) error {
	err := h.Repo.Deny(cmd.Id)
	if err != nil {
		return err
	}

	return nil
}

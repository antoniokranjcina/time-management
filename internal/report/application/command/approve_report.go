package command

import (
	"time-management/internal/report/domain"
)

type ApproveReportCommand struct {
	Id string
}

type ApproveReportHandler struct {
	Repo domain.ReportRepository
}

func (h *ApproveReportHandler) Handle(cmd ApproveReportCommand) error {
	err := h.Repo.Approve(cmd.Id)
	if err != nil {
		return err
	}

	return nil
}

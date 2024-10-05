package command

import (
	"context"
	"time-management/internal/report/domain"
)

type ApproveReportCommand struct {
	Id string
}

type ApproveReportHandler struct {
	Repo domain.ReportRepository
}

func (h *ApproveReportHandler) Handle(ctx context.Context, cmd ApproveReportCommand) error {
	err := h.Repo.Approve(ctx, cmd.Id)
	if err != nil {
		return err
	}

	return nil
}

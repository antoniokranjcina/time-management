package command

import (
	"context"
	"time-management/internal/report/domain"
)

type DenyReportCommand struct {
	Id string
}

type DenyReportHandler struct {
	Repo domain.ReportRepository
}

func (h *DenyReportHandler) Handle(ctx context.Context, cmd DenyReportCommand) error {
	err := h.Repo.Deny(ctx, cmd.Id)
	if err != nil {
		return err
	}

	return nil
}

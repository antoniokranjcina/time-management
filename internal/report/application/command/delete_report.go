package command

import (
	"context"
	"time-management/internal/report/domain"
)

type DeleteReport struct {
	Id string
}

type DeleteReportHandler struct {
	Repo domain.ReportRepository
}

func (h *DeleteReportHandler) Handle(ctx context.Context, cmd DeleteReport) error {
	err := h.Repo.Delete(ctx, cmd.Id)
	if err != nil {
		return err
	}

	return nil
}

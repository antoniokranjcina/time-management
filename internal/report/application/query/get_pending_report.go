package query

import (
	"context"
	"time-management/internal/report/domain"
)

type GetPendingReportQuery struct {
	Id string
}

type GetPendingReportHandler struct {
	Repo domain.ReportRepository
}

func (h *GetPendingReportHandler) Handle(ctx context.Context, query GetPendingReportQuery) (*domain.Report, error) {
	reports, err := h.Repo.GetById(ctx, query.Id, domain.Pending)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

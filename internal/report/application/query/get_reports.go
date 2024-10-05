package query

import (
	"context"
	"time-management/internal/report/domain"
)

type GetReportQuery struct {
	Id string
}

type GetReportHandler struct {
	Repo domain.ReportRepository
}

func (h *GetReportHandler) Handle(ctx context.Context, query GetReportQuery) (*domain.Report, error) {
	report, err := h.Repo.GetById(ctx, query.Id, domain.Approved)
	if err != nil {
		return nil, err
	}

	return report, nil
}

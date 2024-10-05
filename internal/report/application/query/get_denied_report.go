package query

import (
	"context"
	"time-management/internal/report/domain"
)

type GetDeniedReportQuery struct {
	Id string
}

type GetDeniedReportHandler struct {
	Repo domain.ReportRepository
}

func (h *GetDeniedReportHandler) Handle(ctx context.Context, query GetDeniedReportQuery) (*domain.Report, error) {
	report, err := h.Repo.GetById(ctx, query.Id, domain.Denied)
	if err != nil {
		return nil, err
	}

	return report, nil
}

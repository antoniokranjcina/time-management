package query

import (
	"context"
	"time-management/internal/report/domain"
)

type GetDeniedReportsHandler struct {
	Repo domain.ReportRepository
}

func (h *GetDeniedReportsHandler) Handle(ctx context.Context) ([]domain.Report, error) {
	reports, err := h.Repo.GetAll(ctx, domain.Denied)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

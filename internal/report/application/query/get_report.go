package query

import (
	"context"
	"time-management/internal/report/domain"
)

type GetReportsHandler struct {
	Repo domain.ReportRepository
}

func (h *GetReportsHandler) Handle(ctx context.Context) ([]domain.Report, error) {
	reports, err := h.Repo.GetAll(ctx, domain.Approved)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

package query

import (
	"context"
	"time-management/internal/report/domain"
)

type GetPendingReportsHandler struct {
	Repo domain.ReportRepository
}

func (h *GetPendingReportsHandler) Handle(ctx context.Context) ([]domain.Report, error) {
	reports, err := h.Repo.GetAll(ctx, domain.Pending)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

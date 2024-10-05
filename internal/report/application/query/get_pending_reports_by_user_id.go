package query

import (
	"context"
	"time-management/internal/report/domain"
)

type GetPendingReportsByUserIdQuery struct {
	UserId string
}

type GetPendingReportsByUserIdHandler struct {
	Repo domain.ReportRepository
}

func (h *GetPendingReportsByUserIdHandler) Handle(ctx context.Context, query GetPendingReportsByUserIdQuery) ([]domain.Report, error) {
	reports, err := h.Repo.GetAllWithUserId(ctx, query.UserId, domain.Pending)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

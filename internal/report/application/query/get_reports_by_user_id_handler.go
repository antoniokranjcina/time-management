package query

import (
	"context"
	"time-management/internal/report/domain"
)

type GetReportsByUserIdQuery struct {
	UserId string
}

type GetReportsByUserIdHandler struct {
	Repo domain.ReportRepository
}

func (h *GetReportsByUserIdHandler) Handle(
	ctx context.Context,
	query GetReportsByUserIdQuery,
) ([]domain.Report, error) {
	reports, err := h.Repo.GetAllWithUserId(ctx, query.UserId, domain.Approved)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

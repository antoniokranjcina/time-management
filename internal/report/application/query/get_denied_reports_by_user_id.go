package query

import (
	"context"
	"time-management/internal/report/domain"
)

type GetDeniedReportsByUserIdQuery struct {
	UserId string
}

type GetDeniedReportsByUserIdHandler struct {
	Repo domain.ReportRepository
}

func (h *GetDeniedReportsByUserIdHandler) Handle(
	ctx context.Context,
	query GetDeniedReportsByUserIdQuery,
) ([]domain.Report, error) {
	reports, err := h.Repo.GetAllWithUserId(ctx, query.UserId, domain.Denied)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

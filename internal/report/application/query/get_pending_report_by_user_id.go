package query

import (
	"context"
	"time-management/internal/report/domain"
)

type GetPendingReportByUserIdQuery struct {
	Id     string
	UserId string
}

type GetPendingReportByUserIdHandler struct {
	Repo domain.ReportRepository
}

func (h *GetPendingReportByUserIdHandler) Handle(
	ctx context.Context,
	query GetPendingReportByUserIdQuery,
) (*domain.Report, error) {
	report, err := h.Repo.GetByIdWithUserId(ctx, query.Id, query.UserId, domain.Pending)
	if err != nil {
		return nil, err
	}

	return report, nil
}

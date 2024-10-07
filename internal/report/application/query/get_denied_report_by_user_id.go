package query

import (
	"context"
	"time-management/internal/report/domain"
)

type GetDeniedReportByUserIdQuery struct {
	Id     string
	UserId string
}

type GetDeniedReportByUserIdHandler struct {
	Repo domain.ReportRepository
}

func (h *GetDeniedReportByUserIdHandler) Handle(
	ctx context.Context,
	query GetDeniedReportByUserIdQuery,
) (*domain.Report, error) {
	reports, err := h.Repo.GetByIdWithUserId(ctx, query.Id, query.UserId, domain.Denied)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

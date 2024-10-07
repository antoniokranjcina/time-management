package query

import (
	"context"
	"time-management/internal/report/domain"
)

type GetReportByUserIdQuery struct {
	Id     string
	UserId string
}

type GetReportByUserIdHandler struct {
	Repo domain.ReportRepository
}

func (h *GetReportByUserIdHandler) Handle(
	ctx context.Context,
	query GetReportByUserIdQuery,
) (*domain.Report, error) {
	reports, err := h.Repo.GetByIdWithUserId(ctx, query.Id, query.UserId, domain.Approved)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

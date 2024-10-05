package query

import "time-management/internal/report/domain"

type GetPendingReportByUserIdQuery struct {
	Id     string
	UserId string
}

type GetPendingReportByUserIdHandler struct {
	Repo domain.ReportRepository
}

func (h *GetPendingReportByUserIdHandler) Handle(query GetPendingReportByUserIdQuery) (*domain.Report, error) {
	reports, err := h.Repo.GetByIdWithUserId(query.Id, query.UserId, domain.Pending)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

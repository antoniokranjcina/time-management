package query

import "time-management/internal/report/domain"

type GetPendingReportsByUserIdQuery struct {
	UserId string
}

type GetPendingReportsByUserIdHandler struct {
	Repo domain.ReportRepository
}

func (h *GetPendingReportsByUserIdHandler) Handle(query GetPendingReportsByUserIdQuery) ([]domain.Report, error) {
	reports, err := h.Repo.GetAllWithUserId(query.UserId, domain.Pending)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

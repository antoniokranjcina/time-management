package query

import "time-management/internal/report/domain"

type GetPendingReportQuery struct {
	Id string
}

type GetPendingReportHandler struct {
	Repo domain.ReportRepository
}

func (h *GetPendingReportHandler) Handle(query GetPendingReportQuery) (*domain.Report, error) {
	reports, err := h.Repo.GetById(query.Id, domain.Pending)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

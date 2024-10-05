package query

import "time-management/internal/report/domain"

type GetReportQuery struct {
	Id string
}

type GetReportHandler struct {
	Repo domain.ReportRepository
}

func (h *GetReportHandler) Handle(query GetReportQuery) (*domain.Report, error) {
	report, err := h.Repo.GetById(query.Id, domain.Approved)
	if err != nil {
		return nil, err
	}

	return report, nil
}

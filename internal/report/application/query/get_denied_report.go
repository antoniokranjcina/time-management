package query

import "time-management/internal/report/domain"

type GetDeniedReportQuery struct {
	Id string
}

type GetDeniedReportHandler struct {
	Repo domain.ReportRepository
}

func (h *GetDeniedReportHandler) Handle(query GetDeniedReportQuery) (*domain.Report, error) {
	report, err := h.Repo.GetDeniedById(query.Id)
	if err != nil {
		return nil, err
	}

	return report, nil
}

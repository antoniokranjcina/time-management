package query

import "time-management/internal/report/domain"

type GetDeniedReportsHandler struct {
	Repo domain.ReportRepository
}

func (h *GetDeniedReportsHandler) Handle() ([]domain.Report, error) {
	reports, err := h.Repo.GetDeniedAll()
	if err != nil {
		return nil, err
	}

	return reports, nil
}

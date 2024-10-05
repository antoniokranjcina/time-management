package query

import "time-management/internal/report/domain"

type GetDeniedReportsHandler struct {
	Repo domain.ReportRepository
}

func (h *GetDeniedReportsHandler) Handle() ([]domain.Report, error) {
	reports, err := h.Repo.GetAll(domain.Denied)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

package query

import "time-management/internal/report/domain"

type GetReportsHandler struct {
	Repo domain.ReportRepository
}

func (h *GetReportsHandler) Handle() ([]domain.Report, error) {
	reports, err := h.Repo.GetAll(domain.Approved)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

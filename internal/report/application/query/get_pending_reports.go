package query

import "time-management/internal/report/domain"

type GetPendingReportsHandler struct {
	Repo domain.ReportRepository
}

func (h *GetPendingReportsHandler) Handle() ([]domain.Report, error) {
	reports, err := h.Repo.GetAll(domain.Pending)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

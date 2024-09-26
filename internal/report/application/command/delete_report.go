package command

import "time-management/internal/report/domain"

type DeleteReport struct {
	Id string
}

type DeleteReportHandler struct {
	Repo domain.ReportRepository
}

func (h *DeleteReportHandler) Handle(cmd DeleteReport) error {
	err := h.Repo.Delete(cmd.Id)
	if err != nil {
		return err
	}

	return nil
}

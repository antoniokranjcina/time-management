package command

import (
	"context"
	"time-management/internal/report/domain"
	"time-management/internal/shared/util"
)

type UpdatePendingReportCommand struct {
	Id               string
	LocationId       string
	WorkingHours     uint64
	MaintenanceHours uint64
}

type UpdatePendingReportHandler struct {
	Repo domain.ReportRepository
}

func (h *UpdatePendingReportHandler) Handle(ctx context.Context, cmd UpdatePendingReportCommand) (*domain.Report, error) {
	if cmd.LocationId == "" || len(cmd.LocationId) >= 50 {
		return nil, util.NewValidationError(domain.ErrWrongLocationId)
	}
	if cmd.WorkingHours <= 0 || cmd.WorkingHours > 16 {
		return nil, util.NewValidationError(domain.ErrInvalidWorkingHours)
	}
	if cmd.MaintenanceHours < 0 || cmd.MaintenanceHours > 16 {
		return nil, util.NewValidationError(domain.ErrInvalidMaintenanceHours)
	}
	if cmd.WorkingHours+cmd.MaintenanceHours > 16 {
		return nil, util.NewValidationError(domain.ErrInvalidHoursSum)
	}

	updatedReport, err := h.Repo.Update(
		ctx,
		cmd.Id,
		cmd.LocationId,
		cmd.WorkingHours,
		cmd.MaintenanceHours,
		domain.Pending,
	)
	if err != nil {
		return nil, err
	}

	return updatedReport, nil
}

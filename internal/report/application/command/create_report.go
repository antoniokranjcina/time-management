package command

import (
	"context"
	"github.com/google/uuid"
	"time"
	"time-management/internal/report/domain"
	"time-management/internal/shared/util"
)

type CreateReportCommand struct {
	EmployeeId       string
	LocationId       string
	WorkingHours     uint64
	MaintenanceHours uint64
}

type CreateReportHandler struct {
	Repo domain.ReportRepository
}

func (h *CreateReportHandler) Handle(ctx context.Context, cmd CreateReportCommand) (*domain.Report, error) {
	if cmd.EmployeeId == "" || len(cmd.EmployeeId) >= 50 {
		return nil, util.NewValidationError(domain.ErrWrongEmployeeId)
	}
	if cmd.LocationId == "" || len(cmd.LocationId) >= 50 {
		return nil, util.NewValidationError(domain.ErrWrongLocationId)
	}
	if cmd.WorkingHours < 0 || cmd.WorkingHours > 16 {
		return nil, util.NewValidationError(domain.ErrInvalidWorkingHours)
	}
	if cmd.MaintenanceHours <= 0 || cmd.MaintenanceHours > 16 {
		return nil, util.NewValidationError(domain.ErrInvalidMaintenanceHours)
	}
	if cmd.WorkingHours+cmd.MaintenanceHours > 16 {
		return nil, util.NewValidationError(domain.ErrInvalidHoursSum)
	}

	report := domain.NewReport(
		uuid.New().String(),
		cmd.EmployeeId,
		cmd.LocationId,
		cmd.WorkingHours,
		cmd.MaintenanceHours,
		domain.Pending,
		uint64(time.Now().Unix()),
	)

	createdReport, err := h.Repo.Create(ctx, report)
	if err != nil {
		return nil, err
	}

	return createdReport, nil
}

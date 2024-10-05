package domain

import "context"

type ReportRepository interface {
	Create(ctx context.Context, report *Report) (*Report, error)
	GetAll(ctx context.Context, status ReportStatus) ([]Report, error)
	GetAllWithUserId(ctx context.Context, employeeId string, status ReportStatus) ([]Report, error)
	GetById(ctx context.Context, id string, status ReportStatus) (*Report, error)
	GetByIdWithUserId(ctx context.Context, id, userId string, status ReportStatus) (*Report, error)
	Update(
		ctx context.Context,
		id, locationId string,
		workingHours, maintenanceHours uint64,
		status ReportStatus,
	) (*Report, error)
	Approve(ctx context.Context, id string) error
	Deny(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

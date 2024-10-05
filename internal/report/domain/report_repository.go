package domain

type ReportRepository interface {
	Create(report *Report) (*Report, error)
	GetAll(status ReportStatus) ([]Report, error)
	GetAllWithUserId(employeeId string, status ReportStatus) ([]Report, error)
	GetById(id string, status ReportStatus) (*Report, error)
	GetByIdWithUserId(id, userId string, status ReportStatus) (*Report, error)
	Update(id, locationId string, workingHours, maintenanceHours uint64, status ReportStatus) (*Report, error)
	Approve(id string) error
	Deny(id string) error
	Delete(id string) error
}

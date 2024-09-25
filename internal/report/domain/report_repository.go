package domain

type ReportRepository interface {
	Create(report *Report) (*Report, error)
	GetAll() ([]Report, error)
	GetPendingAll() ([]Report, error)
	GetDeniedAll() ([]Report, error)
	GetById(id string) (*Report, error)
	GetPendingById(id string) (*Report, error)
	GetDeniedById(id string) (*Report, error)
	UpdatePending(id, locationId string, workingHours, maintenanceHours uint64) (*Report, error)
	Approve(id string) error
	Deny(id string) error
	Delete(id string) error
}

package domain

type Report struct {
	Id               string       `json:"id"`
	Employee         Employee     `json:"employee"`
	Location         Location     `json:"location"`
	WorkingHours     uint64       `json:"working_hours"`
	MaintenanceHours uint64       `json:"maintenance_hours"`
	Status           ReportStatus `json:"status"`
	CreatedAt        uint64       `json:"created_at"`
}

type Employee struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type Location struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func NewReport(
	id string,
	employeeId string,
	locationId string,
	workingHours uint64,
	maintenanceHours uint64,
	status ReportStatus,
	createdAt uint64,
) *Report {
	return &Report{
		Id:               id,
		Employee:         Employee{Id: employeeId},
		Location:         Location{Id: locationId},
		WorkingHours:     workingHours,
		MaintenanceHours: maintenanceHours,
		Status:           status,
		CreatedAt:        createdAt,
	}
}

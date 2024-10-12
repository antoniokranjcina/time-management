package domain

type Report struct {
	Id               string       `json:"id"`
	User             User         `json:"user"`
	Location         Location     `json:"location"`
	WorkingHours     uint64       `json:"working_hours"`
	MaintenanceHours uint64       `json:"maintenance_hours"`
	Status           ReportStatus `json:"status"`
	CreatedAt        uint64       `json:"created_at"`
}

type User struct {
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
	userId string,
	locationId string,
	workingHours uint64,
	maintenanceHours uint64,
	status ReportStatus,
	createdAt uint64,
) *Report {
	return &Report{
		Id:               id,
		User:             User{Id: userId},
		Location:         Location{Id: locationId},
		WorkingHours:     workingHours,
		MaintenanceHours: maintenanceHours,
		Status:           status,
		CreatedAt:        createdAt,
	}
}

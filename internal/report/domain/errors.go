package domain

import "errors"

var (
	ErrReportNotFound               = errors.New("report not found")
	ErrWrongEmployeeId              = errors.New("wrong employee id: employee does not exist")
	ErrWrongLocationId              = errors.New("wrong location id: location does not exist")
	ErrInvalidWorkingHours          = errors.New("invalid working hours")
	ErrInvalidMaintenanceHours      = errors.New("invalid maintenance hours")
	ErrInvalidHoursInput            = errors.New("invalid hours input")
	ErrInvalidHoursSum              = errors.New("invalid hours sum")
	ErrCannotUpdateReport           = errors.New("cannot update report which is approved or denied")
	ErrReportNotFoundOrUnauthorized = errors.New("report not found")
)

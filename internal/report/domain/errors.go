package domain

import "errors"

var (
	ErrPendingReportNotFound   = errors.New("pending report not found")
	ErrDeniedReportNotFound    = errors.New("denied report not found")
	ErrReportNotFound          = errors.New("report not found")
	ErrWrongEmployeeId         = errors.New("wrong employee id: employee does not exist")
	ErrWrongLocationId         = errors.New("wrong location id: location does not exist")
	ErrInvalidWorkingHours     = errors.New("invalid working hours")
	ErrInvalidMaintenanceHours = errors.New("invalid maintenance hours")
	ErrInvalidHoursSum         = errors.New("invalid hours sum")
	ErrCannotUpdateReport      = errors.New("cannot update report which is approved or denied")
)

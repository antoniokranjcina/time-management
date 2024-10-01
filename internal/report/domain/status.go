package domain

import "fmt"

// ReportStatus defines the possible statuses for a report
type ReportStatus int

const (
	// Pending Enum values for ReportStatus
	Pending  ReportStatus = iota // 0
	Approved                     // 1
	Denied                       // 2
)

// To convert the ReportStatus to a string
func (s ReportStatus) String() string {
	return [...]string{"pending", "approved", "denied"}[s]
}

// ParseReportStatus For parsing a string back to ReportStatus
func ParseReportStatus(status string) (ReportStatus, error) {
	switch status {
	case "pending":
		return Pending, nil
	case "approved":
		return Approved, nil
	case "denied":
		return Denied, nil
	default:
		return -1, fmt.Errorf("invalid report status: %s", status)
	}
}

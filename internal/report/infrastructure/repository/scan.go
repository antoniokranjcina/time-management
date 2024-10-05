package repository

import (
	"database/sql"
	"time-management/internal/report/domain"
)

func (r *PgReportRepository) ScanReportRow(row *sql.Row) (*domain.Report, error) {
	var report domain.Report
	var employee domain.Employee
	var location domain.Location

	err := row.Scan(
		&report.Id, &report.WorkingHours, &report.MaintenanceHours, &report.Status, &report.CreatedAt,
		&employee.Id, &employee.FirstName, &employee.LastName, &employee.Email,
		&location.Id, &location.Name,
	)
	if err != nil {
		return nil, err
	}

	report.Employee = employee
	report.Location = location

	return &report, nil
}

func (r *PgReportRepository) ScanReportRows(rows *sql.Rows) ([]domain.Report, error) {
	var reports []domain.Report

	for rows.Next() {
		var report domain.Report
		var employee domain.Employee
		var location domain.Location

		err := rows.Scan(
			&report.Id, &report.WorkingHours, &report.MaintenanceHours, &report.Status, &report.CreatedAt,
			&employee.Id, &employee.FirstName, &employee.LastName, &employee.Email,
			&location.Id, &location.Name,
		)
		if err != nil {
			return nil, err
		}

		report.Employee = employee
		report.Location = location
		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reports, nil
}

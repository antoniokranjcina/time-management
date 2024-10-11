package repository

import (
	"database/sql"
	"errors"
	"time-management/internal/report/domain"
	"time-management/internal/shared/util"
)

func (r *PgReportRepository) ScanReportRow(row *sql.Row) (*domain.Report, error) {
	var report domain.Report
	var employee domain.User
	var location domain.Location

	err := row.Scan(
		&report.Id, &report.WorkingHours, &report.MaintenanceHours, &report.Status, &report.CreatedAt,
		&employee.Id, &employee.FirstName, &employee.LastName, &employee.Email,
		&location.Id, &location.Name,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, util.NewNotFoundError(domain.ErrReportNotFound)
		}

		return nil, err
	}

	report.User = employee
	report.Location = location

	return &report, nil
}

func (r *PgReportRepository) ScanReportRows(rows *sql.Rows) ([]domain.Report, error) {
	var reports []domain.Report

	for rows.Next() {
		var report domain.Report
		var user domain.User
		var location domain.Location

		err := rows.Scan(
			&report.Id, &report.WorkingHours, &report.MaintenanceHours, &report.Status, &report.CreatedAt,
			&user.Id, &user.FirstName, &user.LastName, &user.Email,
			&location.Id, &location.Name,
		)
		if err != nil {
			return nil, err
		}

		report.User = user
		report.Location = location
		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reports, nil
}

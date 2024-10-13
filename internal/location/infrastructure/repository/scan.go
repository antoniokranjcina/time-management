package repository

import (
	"database/sql"
	"time-management/internal/location/domain"
)

func ScanLocationRow(row *sql.Row) (*domain.Location, error) {
	location := &domain.Location{}
	err := row.Scan(&location.Id, &location.Name, &location.CreatedAt)
	if err != nil {
		return nil, err
	}

	return location, nil
}

func ScanLocationRows(rows *sql.Rows) ([]domain.Location, error) {
	var locations []domain.Location

	for rows.Next() {
		var location domain.Location
		err := rows.Scan(&location.Id, &location.Name, &location.CreatedAt)
		if err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}

	return locations, nil
}

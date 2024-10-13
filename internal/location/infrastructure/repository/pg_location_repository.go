package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time-management/internal/location/domain"
	"time-management/internal/shared/util"
)

const TableName = "locations"

type PgLocationRepository struct {
	DB *sql.DB
}

func NewPgLocationRepository(db *sql.DB) *PgLocationRepository {
	repository := &PgLocationRepository{DB: db}
	err := repository.createLocationTable()
	if err != nil {
		panic(err)
		return nil
	}

	return repository
}

func (r *PgLocationRepository) createLocationTable() error {
	query := fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s (
		id VARCHAR(50) PRIMARY KEY,
		name VARCHAR(50),
		created_at SERIAL
	)`, TableName)

	_, err := r.DB.Exec(query)
	return err
}

func (r *PgLocationRepository) Create(ctx context.Context, location *domain.Location) (*domain.Location, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (id, name, created_at) 
		VALUES ($1, $2, $3) 
		RETURNING id, name, created_at
	`, TableName)

	row := r.DB.QueryRowContext(ctx, query, location.Id, location.Name, location.CreatedAt)
	savedLocation, err := ScanLocationRow(row)
	if err != nil {
		return nil, err
	}

	return savedLocation, nil
}

func (r *PgLocationRepository) GetAll(ctx context.Context) ([]domain.Location, error) {
	query := fmt.Sprintf(`SELECT * FROM %s`, TableName)

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	locations, err := ScanLocationRows(rows)
	if err != nil {
		return nil, err
	}

	return locations, nil
}

func (r *PgLocationRepository) GetById(ctx context.Context, id string) (*domain.Location, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, TableName)

	row := r.DB.QueryRowContext(ctx, query, id)
	location, err := ScanLocationRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, util.NewValidationError(domain.ErrLocationNotFound)
		}
		return nil, err
	}

	return location, nil
}

func (r *PgLocationRepository) Update(ctx context.Context, id, name string) (*domain.Location, error) {
	query := fmt.Sprintf(`UPDATE %s SET name = $1 WHERE id = $2 RETURNING id, name, created_at`, TableName)

	row := r.DB.QueryRowContext(ctx, query, name, id)
	location, err := ScanLocationRow(row)
	if err != nil {
		return nil, err
	}

	return location, nil
}

func (r *PgLocationRepository) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, TableName)

	_, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

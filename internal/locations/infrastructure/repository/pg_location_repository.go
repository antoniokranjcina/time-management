package repository

import (
	"database/sql"
	"errors"
	"log"
	"time-management/internal/locations/domain"
)

type PgLocationRepository struct {
	DB *sql.DB
}

func NewPgLocationRepository(db *sql.DB) *PgLocationRepository {
	return &PgLocationRepository{DB: db}
}

func (r *PgLocationRepository) createLocationTable() error {
	query := `create table if not exists locations (
		id varchar(50) primary key,
		name varchar(50),
		created_at serial
	)`

	_, err := r.DB.Exec(query)
	return err
}

func (r *PgLocationRepository) GetById(id string) (*domain.Location, error) {
	query := `SELECT id, name, created_at FROM locations WHERE id = $1`
	row := r.DB.QueryRow(query, id)

	location := &domain.Location{}
	err := row.Scan(&location.Id, &location.Name, &location.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrLocationNotFound
		}
		return nil, err
	}

	return location, nil
}

func (r *PgLocationRepository) GetAll() ([]domain.Location, error) {
	query := `SELECT * FROM locations`
	rows, err := r.DB.Query(query)
	if err != nil {
		log.Println("Error fetching locations:", err)
		return nil, err
	}
	defer rows.Close()

	var locations []domain.Location
	for rows.Next() {
		var location domain.Location
		if err := rows.Scan(&location.Id, &location.Name, &location.CreatedAt); err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}

	return locations, nil
}

func (r *PgLocationRepository) Save(location *domain.Location) (*domain.Location, error) {
	query := `
		INSERT INTO locations (id, name, created_at) 
		VALUES ($1, $2, $3) 
		RETURNING id, name, created_at
	`
	loc := &domain.Location{}

	err := r.DB.QueryRow(
		query,
		location.Id,
		location.Name,
		location.CreatedAt,
	).Scan(&loc.Id, &loc.Name, &loc.CreatedAt)
	if err != nil {
		log.Println("Error adding location:", err)
		return nil, err
	}

	return loc, nil
}

func (r *PgLocationRepository) Update(id string, name string) (*domain.Location, error) {
	query := `UPDATE locations SET name = $1 WHERE id = $2 RETURNING id, name, created_at`
	loc := &domain.Location{}

	err := r.DB.QueryRow(query, name, id).Scan(&loc.Id, &loc.Name, &loc.CreatedAt)
	if err != nil {
		log.Println("Error updating location:", err)
		return nil, err
	}

	return loc, nil
}

func (r *PgLocationRepository) Delete(id string) error {
	query := `DELETE FROM locations WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		log.Println("Error deleting location:", err)
		return err
	}
	return nil
}

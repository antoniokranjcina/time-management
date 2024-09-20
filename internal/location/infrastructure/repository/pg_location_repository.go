package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time-management/internal/location/domain"
)

const tableName = "locations"

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
	query := fmt.Sprintf(`create table if not exists %s (
		id varchar(50) primary key,
		name varchar(50),
		created_at serial
	)`, tableName)

	_, err := r.DB.Exec(query)
	return err
}

func (r *PgLocationRepository) Save(location *domain.Location) (*domain.Location, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (id, name, created_at) 
		VALUES ($1, $2, $3) 
		RETURNING id, name, created_at
	`, tableName)

	row := r.DB.QueryRow(query, location.Id, location.Name, location.CreatedAt)
	savedLocation, err := scanLocationRow(row)
	if err != nil {
		log.Println("Error adding location:", err)
		return nil, err
	}

	return savedLocation, nil
}

func (r *PgLocationRepository) GetAll() ([]domain.Location, error) {
	query := fmt.Sprintf(`SELECT * FROM %s`, tableName)

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	locations, err := scanLocationRows(rows)
	if err != nil {
		return nil, err
	}

	return locations, nil
}

func (r *PgLocationRepository) GetById(id string) (*domain.Location, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, tableName)

	row := r.DB.QueryRow(query, id)
	location, err := scanLocationRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrLocationNotFound
		}
		return nil, err
	}

	return location, nil
}

func (r *PgLocationRepository) Update(id string, name string) (*domain.Location, error) {
	query := fmt.Sprintf(`UPDATE %s SET name = $1 WHERE id = $2 RETURNING id, name, created_at`, tableName)

	row := r.DB.QueryRow(query, name, id)
	location, err := scanLocationRow(row)
	if err != nil {
		return nil, err
	}

	return location, nil
}

func (r *PgLocationRepository) Delete(id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, tableName)

	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func scanLocationRow(row *sql.Row) (*domain.Location, error) {
	location := &domain.Location{}
	err := row.Scan(&location.Id, &location.Name, &location.CreatedAt)
	if err != nil {
		return nil, err
	}

	return location, nil
}

func scanLocationRows(rows *sql.Rows) ([]domain.Location, error) {
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

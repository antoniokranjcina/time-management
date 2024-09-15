package database

import (
	"github.com/google/uuid"
	"log"
	"time"
)

type Location struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt uint64 `json:"created_at"`
}

func (s *DbService) Init() error {
	return s.createLocationTable()
}

func (s *DbService) createLocationTable() error {
	query := `create table if not exists locations (
		id varchar(50) primary key,
		name varchar(50),
		created_at serial
	)`

	_, err := s.DB.Exec(query)
	return err
}

func (s *DbService) GetLocationById(id string) (*Location, error) {
	query := `SELECT id, name, created_at FROM locations WHERE id = $1`

	location := &Location{}
	err := s.DB.QueryRow(query, id).Scan(&location.Id, &location.Name, &location.CreatedAt)
	if err != nil {
		return nil, err
	}

	return location, nil
}

func (s *DbService) GetLocations() ([]Location, error) {
	query := `SELECT * FROM locations`
	rows, err := s.DB.Query(query)
	if err != nil {
		log.Println("Error fetching locations:", err)
		return nil, err
	}
	defer rows.Close()

	var locations []Location
	for rows.Next() {
		var location Location
		if err := rows.Scan(&location.Id, &location.Name, &location.CreatedAt); err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}

	return locations, nil
}

func (s *DbService) CreateLocation(name string) (*Location, error) {
	query := `
		INSERT INTO locations (id, name, created_at) 
		VALUES ($1, $2, $3) 
		RETURNING id, name, created_at
	`
	loc := &Location{}

	err := s.DB.QueryRow(
		query,
		uuid.New().String(),
		name,
		uint64(time.Now().Unix()),
	).Scan(&loc.Id, &loc.Name, &loc.CreatedAt)
	if err != nil {
		log.Println("Error adding location:", err)
		return nil, err
	}

	return loc, nil
}

func (s *DbService) UpdateLocation(id, name string) (*Location, error) {
	query := `UPDATE locations SET name = $1 WHERE id = $2 RETURNING id, name, created_at`
	location := &Location{}

	err := s.DB.QueryRow(query, name, id).Scan(&location.Id, &location.Name, &location.CreatedAt)
	if err != nil {
		log.Println("Error updating location:", err)
		return nil, err
	}

	return location, nil
}

func (s *DbService) DeleteLocationById(id string) error {
	query := `DELETE FROM locations WHERE id = $1`
	_, err := s.DB.Exec(query, id)
	if err != nil {
		log.Println("Error deleting location:", err)
		return err
	}
	return nil
}

package database

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
)

// StorageService represents a DbService that interacts with a database.
type StorageService interface {
	GetLocations() ([]Location, error)
	GetLocationById(id string) (*Location, error)
	CreateLocation(name string) (*Location, error)
	UpdateLocation(id, name string) (*Location, error)
	DeleteLocationById(id string) error

	Close() error
}

type DbService struct {
	DB *sql.DB
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	schema     = os.Getenv("DB_SCHEMA")
	dbInstance *DbService
)

func NewConnection() (*DbService, error) {
	if dbInstance != nil {
		return dbInstance, nil
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		username,
		password,
		host,
		port,
		database,
		schema,
	)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	dbInstance = &DbService{
		DB: db,
	}

	err = dbInstance.Init()
	if err != nil {
		return nil, err
	}

	return dbInstance, nil
}

func (s *DbService) Close() error {
	log.Printf("Disconnected from database: %s", database)
	return s.DB.Close()
}

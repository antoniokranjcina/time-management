package database

import (
	"context"
	_ "database/sql"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/docker/go-connections/nat"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var dbService *DbService

func mustStartPostgresContainer() (func(context.Context) error, error) {
	var (
		dbName = "database"
		dbPwd  = "password"
		dbUser = "user"
	)

	dbContainer, err := postgres.Run(
		context.Background(),
		"postgres:latest",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPwd),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	database = dbName
	password = dbPwd
	username = dbUser

	if os.Getenv("DB_SCHEMA") == "" {
		os.Setenv("DB_SCHEMA", "public")
	}
	schema = os.Getenv("DB_SCHEMA")

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, err
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}

	host = dbHost
	port = dbPort.Port()

	return dbContainer.Terminate, err
}

func TestMain(m *testing.M) {
	teardown, err := mustStartPostgresContainer()
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	dbService, err = NewConnection()
	if err != nil {
		log.Fatalf("error establishing database connection: %v", err)
	}

	exitCode := m.Run()

	if err := dbService.Close(); err != nil {
		log.Fatalf("error closing database: %v", err)
	}

	if teardown != nil {
		if err := teardown(context.Background()); err != nil {
			log.Fatalf("could not teardown postgres container: %v", err)
		}
	}

	os.Exit(exitCode)
}

func TestDbService_GetLocations(t *testing.T) {
	err := dbService.Init()
	assert.NoError(t, err)

	loc, err := dbService.CreateLocation("Test Location")
	assert.NoError(t, err)

	locations, err := dbService.GetLocations()
	assert.NoError(t, err)
	assert.NotEmpty(t, locations)

	found := false
	for _, location := range locations {
		if location.Id == loc.Id && location.Name == "Test Location" {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestDbService_GetLocationById(t *testing.T) {
	err := dbService.Init()
	assert.NoError(t, err)

	loc, err := dbService.CreateLocation("Test Location")
	assert.NoError(t, err)

	location, err := dbService.GetLocationById(loc.Id)
	assert.NoError(t, err)
	assert.Equal(t, loc.Name, location.Name)
}

func TestDbService_CreateLocation(t *testing.T) {
	err := dbService.Init()
	assert.NoError(t, err)

	loc, err := dbService.CreateLocation("New Test Location")
	assert.NoError(t, err)
	assert.NotNil(t, loc)
	assert.Equal(t, "New Test Location", loc.Name)
}

func TestDbService_UpdateLocation(t *testing.T) {
	err := dbService.Init()
	assert.NoError(t, err)

	loc, err := dbService.CreateLocation("Old Name")
	assert.NoError(t, err)

	updatedLoc, err := dbService.UpdateLocation(loc.Id, "New Name")
	assert.NoError(t, err)
	assert.Equal(t, "New Name", updatedLoc.Name)
	assert.Equal(t, loc.Id, updatedLoc.Id)
}

func TestDbService_DeleteLocationById(t *testing.T) {
	err := dbService.Init()
	assert.NoError(t, err)

	loc, err := dbService.CreateLocation("To Delete")
	assert.NoError(t, err)

	err = dbService.DeleteLocationById(loc.Id)
	assert.NoError(t, err)

	locations, err := dbService.GetLocations()
	assert.NoError(t, err)

	for _, location := range locations {
		assert.NotEqual(t, loc.Id, location.Id)
	}
}

func TestClose(t *testing.T) {
	if dbService.Close() != nil {
		t.Fatalf("expected Close() to return nil")
	}
}

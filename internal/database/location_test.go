package database

import (
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetLocations(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "created_at"}).
		AddRow("1", "Location A", 1234567890).
		AddRow("2", "Location B", 1234567891)

	mock.ExpectQuery(`SELECT \* FROM locations`).WillReturnRows(rows)

	service := DbService{DB: db}

	locations, err := service.GetLocations()
	assert.NoError(t, err)
	assert.Len(t, locations, 2)
	assert.Equal(t, "Location A", locations[0].Name)
	assert.Equal(t, "Location B", locations[1].Name)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetLocationById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "created_at"}).
		AddRow("1", "Location A", 1234567890)

	mock.ExpectQuery(`SELECT id, name, created_at FROM locations WHERE id = \$1`).
		WithArgs("1").
		WillReturnRows(rows)

	service := DbService{DB: db}

	location, err := service.GetLocationById("1")

	assert.NoError(t, err)
	assert.Equal(t, "1", location.Id)
	assert.Equal(t, "Location A", location.Name)
	assert.Equal(t, uint64(1234567890), location.CreatedAt)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCreateLocation(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	name := "Test Location"
	createdAt := uint64(time.Now().Unix())

	mock.ExpectQuery(`INSERT INTO locations \(id, name, created_at\)`).
		WithArgs(sqlmock.AnyArg(), name, createdAt).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at"}).
			AddRow(uuid.New().String(), name, createdAt))

	service := DbService{DB: db}

	location, err := service.CreateLocation(name)
	assert.NoError(t, err)
	assert.Equal(t, name, location.Name)
	assert.Equal(t, createdAt, location.CreatedAt)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDeleteLocationById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec(`DELETE FROM locations WHERE id = \$1`).WithArgs("1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	service := DbService{DB: db}

	err = service.DeleteLocationById("1")
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateLocation(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	name := "Updated Location"
	id := uuid.New().String()
	createdAt := uint64(time.Now().Unix())

	mock.ExpectQuery(`UPDATE locations SET name = \$1 WHERE id = \$2 RETURNING id, name, created_at`).
		WithArgs(name, id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at"}).AddRow(id, name, createdAt))

	service := DbService{DB: db}

	location, err := service.UpdateLocation(id, name)
	assert.NoError(t, err)
	assert.NotNil(t, location)
	assert.Equal(t, id, location.Id)
	assert.Equal(t, name, location.Name)
	assert.Equal(t, createdAt, location.CreatedAt)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestInit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec(`create table if not exists locations`).WillReturnResult(sqlmock.NewResult(0, 0))

	service := DbService{DB: db}

	err = service.Init()
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

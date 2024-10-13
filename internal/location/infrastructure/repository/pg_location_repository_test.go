package repository

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time-management/internal/location/domain"
)

func TestPgLocationRepository_Create(t *testing.T) {
	mock, repo := setupMockAndRepo(t)

	query := fmt.Sprintf(`INSERT INTO %s`, TableName)
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(loc.Id, loc.Name, loc.CreatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at"}).AddRow(loc.Id, loc.Name, loc.CreatedAt))

	// Execute test
	ctx := context.Background()
	createdLocation, err := repo.Create(ctx, &loc)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, loc.Id, createdLocation.Id)
	assertMockExpectations(t, mock)
}

func TestPgLocationRepository_GetAll(t *testing.T) {
	mock, repo := setupMockAndRepo(t)

	query := fmt.Sprintf(`SELECT * FROM %s`, TableName)
	rows := sqlmock.NewRows([]string{"id", "name", "created_at"})

	locations := []domain.Location{loc}
	for _, loc := range locations {
		rows.AddRow(loc.Id, loc.Name, loc.CreatedAt)
	}
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

	// Execute test
	ctx := context.Background()
	fetchedLocations, err := repo.GetAll(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, fetchedLocations, 1)
	assert.Equal(t, fetchedLocations[0].Id, locations[0].Id)
	assert.Equal(t, fetchedLocations[0].Name, locations[0].Name)
	assert.Equal(t, fetchedLocations[0].CreatedAt, locations[0].CreatedAt)
	assertMockExpectations(t, mock)
}

func TestPgLocationRepository_GetById(t *testing.T) {
	mock, repo := setupMockAndRepo(t)

	query := fmt.Sprintf(`SELECT * FROM %s WHERE id =`, TableName)
	rows := sqlmock.NewRows([]string{"id", "name", "created_at"}).AddRow(loc.Id, loc.Name, loc.CreatedAt)
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(loc.Id).WillReturnRows(rows)

	// Execute test
	ctx := context.Background()
	fetchedLocation, err := repo.GetById(ctx, loc.Id)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, loc.Id, fetchedLocation.Id)
	assert.Equal(t, loc.Name, fetchedLocation.Name)
	assert.Equal(t, loc.CreatedAt, fetchedLocation.CreatedAt)
	assertMockExpectations(t, mock)
}

func TestPgLocationRepository_Update(t *testing.T) {
	mock, repo := setupMockAndRepo(t)

	query := fmt.Sprintf(`UPDATE %s SET name = $1 WHERE id = $2 RETURNING id, name, created_at`, TableName)
	rows := sqlmock.NewRows([]string{"id", "name", "created_at"}).AddRow(loc.Id, loc.Name, loc.CreatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(loc.Name, loc.Id).WillReturnRows(rows)

	// Execute test
	ctx := context.Background()
	fetchedLocation, err := repo.Update(ctx, loc.Id, loc.Name)

	// Assertion
	assert.NoError(t, err)
	assert.Equal(t, loc.Id, fetchedLocation.Id)
	assert.Equal(t, loc.Name, fetchedLocation.Name)
	assert.Equal(t, loc.CreatedAt, fetchedLocation.CreatedAt)
	assertMockExpectations(t, mock)
}

func TestPgLocationRepository_Delete(t *testing.T) {
	mock, repo := setupMockAndRepo(t)

	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, TableName)

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(loc.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute test
	ctx := context.Background()
	err := repo.Delete(ctx, loc.Id)

	// Assertions
	assert.NoError(t, err)
	assertMockExpectations(t, mock)
}

func setupMockAndRepo(t *testing.T) (sqlmock.Sqlmock, *PgLocationRepository) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS locations").
		WillReturnResult(sqlmock.NewResult(0, 0))

	repo := NewPgLocationRepository(db)
	return mock, repo
}

// assertMockExpectations is a helper to ensure all expectations of the mock are met
func assertMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

var loc = domain.Location{
	Id:        "loc123",
	Name:      "New York",
	CreatedAt: 123456789,
}

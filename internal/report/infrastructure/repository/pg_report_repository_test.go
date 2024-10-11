package repository_test

import (
	"context"
	_ "errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	locationPg "time-management/internal/location/infrastructure/repository"
	"time-management/internal/report/domain"
	"time-management/internal/report/infrastructure/repository"
	userPg "time-management/internal/user/infrastructure/repository"
)

func TestPgReportRepository_Create(t *testing.T) {
	mock, repo := setupMockAndRepo(t)

	mockCheckRecordExists(mock, userPg.TableName, rep1.User.Id, true)
	mockCheckRecordExists(mock, locationPg.TableName, rep1.Location.Id, true)

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO reports").
		WithArgs(
			rep1.Id,
			rep1.User.Id,
			rep1.Location.Id,
			rep1.WorkingHours,
			rep1.MaintenanceHours,
			rep1.Status,
			rep1.CreatedAt,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(rep1.Id))
	mock.ExpectCommit()

	mockFullReportQuery(mock, rep1.Id, nil, rep1)

	// Execute test
	ctx := context.Background()
	createdReport, err := repo.Create(ctx, &rep1)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, rep1.Id, createdReport.Id)
	assertMockExpectations(t, mock)
}

func TestPgReportRepository_GetAll(t *testing.T) {
	mock, repo := setupMockAndRepo(t)

	// Mock query response
	mockReportsQuery(mock, nil, domain.Pending, []domain.Report{rep1})

	// Execute test
	ctx := context.Background()
	reports, err := repo.GetAll(ctx, domain.Pending)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, reports, 1)
	assertReportEqual(t, reports[0], rep1.Id, rep1.User.Id, rep1.Location.Id, rep1.WorkingHours, rep1.MaintenanceHours)
	assertMockExpectations(t, mock)
}

func TestPgReportRepository_GetAllWithUserId(t *testing.T) {
	mock, repo := setupMockAndRepo(t)

	userId := "usr123"

	// Mock query response
	mockReportsQuery(mock, &userId, domain.Pending, []domain.Report{rep1, rep2})

	// Execute test
	ctx := context.Background()
	reports, err := repo.GetAllWithUserId(ctx, userId, domain.Pending)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, reports, 2)
	assertReportEqual(t, reports[0], rep1.Id, rep1.User.Id, rep1.Location.Id, rep1.WorkingHours, rep1.MaintenanceHours)
	assertReportEqual(t, reports[1], rep2.Id, rep2.User.Id, rep2.Location.Id, rep2.WorkingHours, rep2.MaintenanceHours)
	assertMockExpectations(t, mock)
}

func TestPgReportRepository_GetById(t *testing.T) {
	mock, repo := setupMockAndRepo(t)

	// Mock query response
	mockReportQuery(mock, rep1.Id, nil, domain.Pending, rep1)

	// Execute test
	ctx := context.Background()
	report, err := repo.GetById(ctx, rep1.Id, domain.Pending)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, rep1.Id, report.Id)
	assertMockExpectations(t, mock)
}

func TestPgReportRepository_GetByIdWithUserId(t *testing.T) {
	mock, repo := setupMockAndRepo(t)

	userId := "user123"

	// Mock query response
	mockReportQuery(mock, rep1.Id, &userId, domain.Approved, rep1)

	// Execute test
	ctx := context.Background()
	report, err := repo.GetByIdWithUserId(ctx, rep1.Id, rep1.User.Id, domain.Approved)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, rep1.Id, report.Id)
	assert.Equal(t, rep1.User.Id, report.User.Id)
	assert.Equal(t, rep1.Location.Name, report.Location.Name)
	assertMockExpectations(t, mock)
}

func TestPgReportRepository_Update(t *testing.T) {
	mock, repo := setupMockAndRepo(t)

	// Input variables
	ctx := context.Background()
	reportId := "report123"
	userId := "user123"
	locationId := "loc123"
	workingHours := uint64(40)
	maintenanceHours := uint64(5)
	status := domain.Pending

	// Mock check location existence
	mockCheckRecordExists(mock, "locations", locationId, true)

	// Mock update query
	mock.ExpectExec(regexp.QuoteMeta(`
        UPDATE reports SET working_hours=$1, maintenance_hours=$2, location_id=$3 
        WHERE id=$4 AND user_id=$5 AND status=$6`)).
		WithArgs(workingHours, maintenanceHours, locationId, reportId, userId, status).
		WillReturnResult(sqlmock.NewResult(1, 1))

	report := domain.Report{
		Id:               reportId,
		User:             domain.User{Id: userId, FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"},
		Location:         domain.Location{Id: locationId, Name: "Location A"},
		WorkingHours:     workingHours,
		MaintenanceHours: maintenanceHours,
		Status:           status,
	}

	// Mock full report query after update
	mockFullReportQuery(mock, reportId, &status, report)

	// Execute test
	updatedReport, err := repo.Update(ctx, reportId, userId, locationId, workingHours, maintenanceHours, status)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, updatedReport)
	assertMockExpectations(t, mock)
}

func TestPgReportRepository_Approve(t *testing.T) {
	mock, repo := setupMockAndRepo(t)

	// Input variables
	ctx := context.Background()
	reportId := "report123"

	// Mock approval query
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE reports SET status = $1 WHERE id = $2`)).
		WithArgs(domain.Approved, reportId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute test
	err := repo.Approve(ctx, reportId)

	// Assertions
	assert.NoError(t, err)
	assertMockExpectations(t, mock)
}

func TestPgReportRepository_Deny(t *testing.T) {
	mock, repo := setupMockAndRepo(t)

	// Input variables
	ctx := context.Background()
	reportId := "report123"

	// Mock denial query
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE reports SET status = $1 WHERE id = $2`)).
		WithArgs(domain.Denied, reportId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute test
	err := repo.Deny(ctx, reportId)

	// Assertions
	assert.NoError(t, err)
	assertMockExpectations(t, mock)
}

func TestPgReportRepository_Delete(t *testing.T) {
	mock, repo := setupMockAndRepo(t)

	// Input variables
	ctx := context.Background()
	reportId := "report123"

	// Mock delete query
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM reports WHERE id = $1`)).
		WithArgs(reportId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute test
	err := repo.Delete(ctx, reportId)

	// Assertions
	assert.NoError(t, err)
	assertMockExpectations(t, mock)
}

func setupMockAndRepo(t *testing.T) (sqlmock.Sqlmock, *repository.PgReportRepository) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS reports").
		WillReturnResult(sqlmock.NewResult(0, 0))

	repo := repository.NewPgReportRepository(db)
	return mock, repo
}

// mockCheckRecordExists mocks the query that checks if a record exists in a given table
func mockCheckRecordExists(mock sqlmock.Sqlmock, tableName, id string, exists bool) {
	query := fmt.Sprintf(`SELECT EXISTS\(SELECT 1 FROM %s WHERE id\s*=\s*\$1\)`, tableName)
	rows := sqlmock.NewRows([]string{"exists"})
	if exists {
		rows.AddRow(true)
	} else {
		rows.AddRow(false)
	}
	mock.ExpectQuery(query).WithArgs(id).WillReturnRows(rows)
}

// mockFullReportQuery mocks the query that retrieves a full report by ID
func mockFullReportQuery(
	mock sqlmock.Sqlmock,
	reportId string,
	status *domain.ReportStatus,
	report domain.Report,
) {
	query :=
		`SELECT 
		    r.id, r.working_hours, r.maintenance_hours, r.status, r.created_at,
			u.id, u.first_name, u.last_name, u.email,
			l.id, l.name
		FROM reports r
		JOIN users u ON r.user_id = u.id
		JOIN locations l ON r.location_id = l.id
		WHERE r.id = $1`

	if status != nil {
		query += " AND r.status = $2"
	}

	// Updated the row columns to match the actual query without column aliases
	rows := sqlmock.NewRows([]string{
		"id", "working_hours", "maintenance_hours", "status", "created_at",
		"id", "first_name", "last_name", "email",
		"id", "name",
	}).AddRow(
		report.Id, report.WorkingHours, report.MaintenanceHours, report.Status, report.CreatedAt,
		report.User.Id, report.User.FirstName, report.User.LastName, report.User.Email,
		report.Location.Id, report.Location.Name,
	)

	if status != nil {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(reportId, &status).WillReturnRows(rows)
	} else {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(reportId).WillReturnRows(rows)
	}
}

// mockReportsQuery mocks the query that retrieves multiple reports by status
func mockReportsQuery(
	mock sqlmock.Sqlmock,
	userId *string,
	status domain.ReportStatus,
	reports []domain.Report,
) {
	query := `
		SELECT 
		    r.id, r.working_hours, r.maintenance_hours, r.status, r.created_at,
			u.id, u.first_name, u.last_name, u.email,
			l.id, l.name
		FROM reports r
		JOIN users u ON r.user_id = u.id
		JOIN locations l ON r.location_id = l.id
		WHERE r.status = $1`

	// Append userId only if it's provided
	if userId != nil {
		query += " AND r.user_id = $2"
	}

	rows := sqlmock.NewRows([]string{
		"id", "working_hours", "maintenance_hours", "status", "created_at",
		"id", "first_name", "last_name", "email",
		"id", "name",
	})

	for _, r := range reports {
		rows.AddRow(
			r.Id, r.WorkingHours, r.MaintenanceHours, r.Status, r.CreatedAt,
			r.User.Id, r.User.FirstName, r.User.LastName, r.User.Email,
			r.Location.Id, r.Location.Name,
		)
	}

	// Conditionally expect query based on whether status is provided
	if userId != nil {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(status, &userId).WillReturnRows(rows)
	} else {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(status).WillReturnRows(rows)
	}
}

// mockReportQuery mocks the query that retrieves a single report by ID and status
func mockReportQuery(
	mock sqlmock.Sqlmock,
	reportId string,
	userId *string,
	status domain.ReportStatus,
	report domain.Report,
) {
	query := `
		SELECT 
		    r.id, r.working_hours, r.maintenance_hours, r.status, r.created_at,
			u.id, u.first_name, u.last_name, u.email,
			l.id, l.name
		FROM reports r
		JOIN users u ON r.user_id = u.id
		JOIN locations l ON r.location_id = l.id
		WHERE r.id = $1 AND r.status = $2`

	// Append userId only if it's provided
	if userId != nil {
		query += " AND r.user_id = $3"
	}

	rows := sqlmock.NewRows([]string{
		"id", "working_hours", "maintenance_hours", "status", "created_at",
		"id", "first_name", "last_name", "email",
		"id", "name",
	}).AddRow(
		report.Id, report.WorkingHours, report.MaintenanceHours, report.Status, report.CreatedAt,
		report.User.Id, report.User.FirstName, report.User.LastName, report.User.Email,
		report.Location.Id, report.Location.Name,
	)

	if userId != nil {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(reportId, status, &userId).WillReturnRows(rows)
	} else {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(reportId, status).WillReturnRows(rows)
	}
}

// assertReportEqual is a helper function for comparing report fields in the assertions
func assertReportEqual(t *testing.T, report domain.Report, id, userId, locationId string, workingHours, maintenanceHours uint64) {
	assert.Equal(t, id, report.Id)
	assert.Equal(t, userId, report.User.Id)
	assert.Equal(t, locationId, report.Location.Id)
	assert.Equal(t, workingHours, report.WorkingHours)
	assert.Equal(t, maintenanceHours, report.MaintenanceHours)
}

// assertMockExpectations is a helper to ensure all expectations of the mock are met
func assertMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

var rep1 = domain.Report{
	Id:               "123",
	User:             domain.User{Id: "user123", FirstName: "John", LastName: "Doe", Email: "john@example.com"},
	Location:         domain.Location{Id: "loc123", Name: "Main Office"},
	WorkingHours:     7,
	MaintenanceHours: 4,
	Status:           domain.Pending,
	CreatedAt:        123456789,
}

var rep2 = domain.Report{
	Id:               "123",
	User:             domain.User{Id: "user123", FirstName: "John", LastName: "Doe", Email: "john@example.com"},
	Location:         domain.Location{Id: "loc123", Name: "Remote"},
	WorkingHours:     8,
	MaintenanceHours: 3,
	Status:           domain.Pending,
	CreatedAt:        123456789,
}

package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	locationPg "time-management/internal/location/infrastructure/repository"
	"time-management/internal/report/domain"
	"time-management/internal/shared/util"
	userPg "time-management/internal/user/infrastructure/repository"
)

const tableName = "reports"

type PgReportRepository struct {
	DB *sql.DB
}

func NewPgReportRepository(db *sql.DB) *PgReportRepository {
	repository := &PgReportRepository{DB: db}
	err := repository.createLocationTable()
	if err != nil {
		panic(err)
		return nil
	}

	return repository
}

func (r *PgReportRepository) createLocationTable() error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id VARCHAR(50) PRIMARY KEY,
			employee_id VARCHAR(50) REFERENCES %s(id) ON DELETE CASCADE,
			location_id VARCHAR(50) REFERENCES %s(id) ON DELETE CASCADE,
			working_hours SERIAL,
			maintenance_hours SERIAL,
			status SERIAL,
			created_at SERIAL
		)`, tableName, userPg.TableName, locationPg.TableName)

	_, err := r.DB.Exec(query)
	return err
}

func (r *PgReportRepository) Create(ctx context.Context, report *domain.Report) (*domain.Report, error) {
	var wg sync.WaitGroup
	wg.Add(2)

	employeeExistChan := make(chan bool)
	locationExistChan := make(chan bool)
	errorChan := make(chan error, 2)

	// Concurrent check for employee and location existence
	go func() {
		defer wg.Done()
		exists, err := r.checkIfRecordExists(ctx, report.Employee.Id, userPg.TableName)
		if err != nil {
			errorChan <- err
			return
		}
		employeeExistChan <- exists
	}()

	go func() {
		defer wg.Done()
		exists, err := r.checkIfRecordExists(ctx, report.Location.Id, locationPg.TableName)
		if err != nil {
			errorChan <- err
			return
		}
		locationExistChan <- exists
	}()

	// Wait for both go routines to finish
	wg.Wait()
	close(employeeExistChan)
	close(locationExistChan)
	close(errorChan)

	var employeeExist, locationExist bool
	for err := range errorChan {
		if err != nil {
			return nil, err
		}
	}
	employeeExist = <-employeeExistChan
	locationExist = <-locationExistChan

	if !employeeExist {
		return nil, util.NewValidationError(domain.ErrWrongEmployeeId)
	}
	if !locationExist {
		return nil, util.NewValidationError(domain.ErrWrongLocationId)
	}

	// Transaction to ensure atomicity
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (id, employee_id, location_id, working_hours, maintenance_hours, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, tableName)

	row := tx.QueryRowContext(
		ctx,
		query, report.Id,
		report.Employee.Id,
		report.Location.Id,
		report.WorkingHours,
		report.MaintenanceHours,
		report.Status,
		report.CreatedAt,
	)

	var savedId string
	if err := row.Scan(&savedId); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.getFullReport(ctx, savedId, nil)
}

func (r *PgReportRepository) GetAll(ctx context.Context, status domain.ReportStatus) ([]domain.Report, error) {
	query := fmt.Sprintf(`
		SELECT 
			r.id, r.working_hours, r.maintenance_hours, r.status, r.created_at,
			e.id, e.first_name, e.last_name, e.email,
			l.id, l.name
		FROM %s r
		JOIN %s e ON r.employee_id = e.id
		JOIN %s l ON r.location_id = l.id
		WHERE r.status = $1;
	`, tableName, userPg.TableName, locationPg.TableName)

	rows, err := r.DB.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.ScanReportRows(rows)
}

func (r *PgReportRepository) GetAllWithUserId(
	ctx context.Context,
	employeeId string,
	status domain.ReportStatus,
) ([]domain.Report, error) {
	query := fmt.Sprintf(`
		SELECT 
			r.id, r.working_hours, r.maintenance_hours, r.status, r.created_at,
			e.id, e.first_name, e.last_name, e.email,
			l.id, l.name
		FROM %s r
		JOIN %s e ON r.employee_id = e.id
		JOIN %s l ON r.location_id = l.id
		WHERE r.employee_id = $1 AND r.status = $2;
	`, tableName, userPg.TableName, locationPg.TableName)

	rows, err := r.DB.QueryContext(ctx, query, employeeId, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.ScanReportRows(rows)
}

func (r *PgReportRepository) GetById(
	ctx context.Context,
	id string,
	status domain.ReportStatus,
) (*domain.Report, error) {
	report, err := r.getFullReportByIdAndStatus(ctx, id, status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, util.NewValidationError(domain.ErrReportNotFound)
		}
	}

	return report, nil
}

func (r *PgReportRepository) GetByIdWithUserId(
	ctx context.Context, id,
	userId string,
	status domain.ReportStatus,
) (*domain.Report, error) {
	query := fmt.Sprintf(`
		SELECT 
			r.id, r.working_hours, r.maintenance_hours, r.status, r.created_at,
			e.id, e.first_name, e.last_name, e.email,
			l.id, l.name
		FROM %s r
		JOIN %s e ON r.employee_id = e.id
		JOIN %s l ON r.location_id = l.id
		WHERE r.employee_id = $1 AND r.id = $2 AND r.status = $3;

	`, tableName, userPg.TableName, locationPg.TableName)

	row := r.DB.QueryRowContext(ctx, query, userId, id, status)

	return r.ScanReportRow(row)
}

func (r *PgReportRepository) Update(
	ctx context.Context,
	id, locationId string,
	workingHours, maintenanceHours uint64,
	status domain.ReportStatus,
) (*domain.Report, error) {
	locationExist, err := r.checkIfRecordExists(ctx, locationId, locationPg.TableName)
	if err != nil {
		return nil, err
	}
	if !locationExist {
		return nil, util.NewValidationError(domain.ErrWrongLocationId)
	}

	query := fmt.Sprintf(`
		UPDATE %s SET working_hours=$1, maintenance_hours=$2, location_id=$3 
	  	WHERE id=$4
		RETURNING id
	`, tableName)

	var updatedId string
	row := r.DB.QueryRowContext(ctx, query, workingHours, maintenanceHours, locationId, id)
	err = row.Scan(&updatedId)
	if err != nil {
		return nil, err
	}

	updateReport, err := r.getFullReportByIdAndStatus(ctx, updatedId, status)
	if err != nil {
		return nil, err
	}

	return updateReport, nil
}

func (r *PgReportRepository) Approve(ctx context.Context, id string) error {
	query := fmt.Sprintf(`UPDATE %s SET status = $1 WHERE id = $2`, tableName)

	_, err := r.DB.ExecContext(ctx, query, domain.Approved, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgReportRepository) Deny(ctx context.Context, id string) error {
	query := fmt.Sprintf(`UPDATE %s SET status = $1 WHERE id = $2`, tableName)

	_, err := r.DB.ExecContext(ctx, query, domain.Denied, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgReportRepository) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, tableName)

	_, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgReportRepository) checkIfRecordExists(ctx context.Context, id, table string) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)`, table)

	var exists bool
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *PgReportRepository) checkIfReportMatchesStatus(
	ctx context.Context,
	id string,
	status domain.ReportStatus,
) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1 AND status = $2)`, tableName)

	var matches bool
	err := r.DB.QueryRowContext(ctx, query, id, status).Scan(&matches)
	if err != nil {
		return false, err
	}

	return matches, nil
}

func (r *PgReportRepository) getFullReportByIdAndStatus(
	ctx context.Context,
	id string,
	status domain.ReportStatus,
) (*domain.Report, error) {
	return r.getFullReport(ctx, id, &status)
}

func (r *PgReportRepository) getFullReport(
	ctx context.Context,
	id string,
	status *domain.ReportStatus,
) (*domain.Report, error) {
	baseQuery := fmt.Sprintf(`
		SELECT 
			r.id, r.working_hours, r.maintenance_hours, r.status, r.created_at,
			e.id, e.first_name, e.last_name, e.email,
			l.id, l.name
		FROM %s r
		JOIN %s e ON r.employee_id = e.id
		JOIN %s l ON r.location_id = l.id
		WHERE r.id = $1
	`, tableName, userPg.TableName, locationPg.TableName)

	if status != nil {
		baseQuery += " AND r.status = $2"
	}

	var row *sql.Row
	if status != nil {
		row = r.DB.QueryRowContext(ctx, baseQuery, id, *status)
	} else {
		row = r.DB.QueryRowContext(ctx, baseQuery, id)
	}

	return r.ScanReportRow(row)
}

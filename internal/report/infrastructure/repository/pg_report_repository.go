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

const TableName = "reports"

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
			user_id VARCHAR(50) REFERENCES %s(id) ON DELETE CASCADE,
			location_id VARCHAR(50) REFERENCES %s(id) ON DELETE CASCADE,
			working_hours SERIAL,
			maintenance_hours SERIAL,
			status SERIAL,
			created_at SERIAL
		)`, TableName, userPg.TableName, locationPg.TableName)

	_, err := r.DB.Exec(query)
	return err
}

func (r *PgReportRepository) Create(ctx context.Context, report *domain.Report) (*domain.Report, error) {
	var wg sync.WaitGroup
	wg.Add(2)

	employeeExistChan := make(chan bool, 1)
	locationExistChan := make(chan bool, 1)
	errorChan := make(chan error, 2)

	// Concurrent check for employee and location existence
	go func() {
		defer wg.Done()
		exists, err := r.checkIfRecordExists(ctx, report.User.Id, userPg.TableName)
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
	go func() {
		wg.Wait()
		close(employeeExistChan)
		close(locationExistChan)
		close(errorChan)
	}()

	// Handle the results from the channels
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
		INSERT INTO %s (id, user_id, location_id, working_hours, maintenance_hours, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, TableName)

	row := tx.QueryRowContext(
		ctx,
		query, report.Id,
		report.User.Id,
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
			u.id, u.first_name, u.last_name, u.email,
			l.id, l.name
		FROM %s r
		JOIN %s u ON r.user_id = u.id
		JOIN %s l ON r.location_id = l.id
		WHERE r.status = $1;
	`, TableName, userPg.TableName, locationPg.TableName)

	rows, err := r.DB.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.ScanReportRows(rows)
}

func (r *PgReportRepository) GetAllWithUserId(
	ctx context.Context,
	userId string,
	status domain.ReportStatus,
) ([]domain.Report, error) {
	query := fmt.Sprintf(`
		SELECT 
			r.id, r.working_hours, r.maintenance_hours, r.status, r.created_at,
			u.id, u.first_name, u.last_name, u.email,
			l.id, l.name
		FROM %s r
		JOIN %s u ON r.user_id = u.id
		JOIN %s l ON r.location_id = l.id
		WHERE r.status = $1 AND r.user_id = $2;
	`, TableName, userPg.TableName, locationPg.TableName)

	rows, err := r.DB.QueryContext(ctx, query, status, userId)
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
	report, err := r.getFullReport(ctx, id, &status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, util.NewNotFoundError(domain.ErrReportNotFound)
		}

		return nil, err
	}

	return report, nil
}

func (r *PgReportRepository) GetByIdWithUserId(
	ctx context.Context,
	id, userId string,
	status domain.ReportStatus,
) (*domain.Report, error) {
	query := fmt.Sprintf(`
		SELECT 
			r.id, r.working_hours, r.maintenance_hours, r.status, r.created_at,
			u.id, u.first_name, u.last_name, u.email,
			l.id, l.name
		FROM %s r
		JOIN %s u ON r.user_id = u.id
		JOIN %s l ON r.location_id = l.id
		WHERE r.id = $1 AND r.status = $2 AND r.user_id = $3;

	`, TableName, userPg.TableName, locationPg.TableName)

	row := r.DB.QueryRowContext(ctx, query, id, status, userId)

	return r.ScanReportRow(row)
}

func (r *PgReportRepository) Update(
	ctx context.Context,
	id, userId, locationId string,
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
		WHERE id=$4 AND user_id=$5 AND status=$6
	`, TableName)

	result, err := r.DB.ExecContext(ctx, query, workingHours, maintenanceHours, locationId, id, userId, status)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, util.NewValidationError(domain.ErrReportNotFoundOrUnauthorized)
	}

	updateReport, err := r.getFullReport(ctx, id, &status)
	if err != nil {
		return nil, err
	}

	return updateReport, nil
}

func (r *PgReportRepository) Approve(ctx context.Context, id string) error {
	query := fmt.Sprintf(`UPDATE %s SET status = $1 WHERE id = $2`, TableName)

	_, err := r.DB.ExecContext(ctx, query, domain.Approved, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgReportRepository) Deny(ctx context.Context, id string) error {
	query := fmt.Sprintf(`UPDATE %s SET status = $1 WHERE id = $2`, TableName)

	_, err := r.DB.ExecContext(ctx, query, domain.Denied, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgReportRepository) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, TableName)

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
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1 AND status = $2)`, TableName)

	var matches bool
	err := r.DB.QueryRowContext(ctx, query, id, status).Scan(&matches)
	if err != nil {
		return false, err
	}

	return matches, nil
}

func (r *PgReportRepository) getFullReport(
	ctx context.Context,
	id string,
	status *domain.ReportStatus,
) (*domain.Report, error) {
	baseQuery := fmt.Sprintf(`
		SELECT 
			r.id, r.working_hours, r.maintenance_hours, r.status, r.created_at,
			u.id, u.first_name, u.last_name, u.email,
			l.id, l.name
		FROM %s r
		JOIN %s u ON r.user_id = u.id
		JOIN %s l ON r.location_id = l.id
		WHERE r.id = $1
	`, TableName, userPg.TableName, locationPg.TableName)

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

package repository

import (
	"database/sql"
	"errors"
	"fmt"
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
			working_hours serial,
			maintenance_hours serial,
			status serial,
			created_at serial
		)`, tableName, userPg.TableName, locationPg.TableName)

	_, err := r.DB.Exec(query)
	return err
}

func (r *PgReportRepository) Create(report *domain.Report) (*domain.Report, error) {
	employeeExistChan := make(chan bool)
	locationExistChan := make(chan bool)
	errorChan := make(chan error)

	go func() {
		exists, err := r.checkIfRecordExists(report.Employee.Id, userPg.TableName)
		if err != nil {
			errorChan <- err
			return
		}
		employeeExistChan <- exists
	}()

	go func() {
		exists, err := r.checkIfRecordExists(report.Location.Id, locationPg.TableName)
		if err != nil {
			errorChan <- err
			return
		}
		locationExistChan <- exists
	}()

	var employeeExist, locationExist bool
	for i := 0; i < 2; i++ {
		select {
		case employeeExist = <-employeeExistChan:
		case locationExist = <-locationExistChan:
		case err := <-errorChan:
			return nil, err
		}
	}

	if !employeeExist {
		return nil, util.NewValidationError(domain.ErrWrongEmployeeId)
	}
	if !locationExist {
		return nil, util.NewValidationError(domain.ErrWrongLocationId)
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (id, employee_id, location_id, working_hours, maintenance_hours, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, tableName)

	row := r.DB.QueryRow(
		query,
		report.Id,
		report.Employee.Id,
		report.Location.Id,
		report.WorkingHours,
		report.MaintenanceHours,
		report.Status,
		report.CreatedAt,
	)

	var savedId string
	err := row.Scan(&savedId)
	if err != nil {
		return nil, err
	}

	rep, err := r.getFullReportById(savedId)
	if err != nil {
		return nil, err
	}

	return rep, nil
}

func (r *PgReportRepository) GetAll(status domain.ReportStatus) ([]domain.Report, error) {
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

	rows, err := r.DB.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.ScanReportRows(rows)
}

func (r *PgReportRepository) GetAllWithUserId(employeeId string, status domain.ReportStatus) ([]domain.Report, error) {
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

	rows, err := r.DB.Query(query, employeeId, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.ScanReportRows(rows)
}

func (r *PgReportRepository) GetById(id string, status domain.ReportStatus) (*domain.Report, error) {
	report, err := r.getFullReportByIdAndStatus(id, status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, util.NewValidationError(domain.ErrReportNotFound)
		}
	}

	return report, nil
}

func (r *PgReportRepository) GetByIdWithUserId(id, userId string, status domain.ReportStatus) (*domain.Report, error) {
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

	row := r.DB.QueryRow(query, userId, id, status)

	return r.ScanReportRow(row)
}

func (r *PgReportRepository) Update(id, locationId string, workingHours, maintenanceHours uint64, status domain.ReportStatus) (*domain.Report, error) {
	locationExist, err := r.checkIfRecordExists(locationId, locationPg.TableName)
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
	row := r.DB.QueryRow(query, workingHours, maintenanceHours, locationId, id)
	err = row.Scan(&updatedId)
	if err != nil {
		return nil, err
	}

	updateReport, err := r.getFullReportByIdAndStatus(updatedId, status)
	if err != nil {
		return nil, err
	}

	return updateReport, nil
}

func (r *PgReportRepository) Approve(id string) error {
	query := fmt.Sprintf(`UPDATE %s SET status = $1 WHERE id = $2`, tableName)

	_, err := r.DB.Exec(query, domain.Approved, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgReportRepository) Deny(id string) error {
	query := fmt.Sprintf(`UPDATE %s SET status = $1 WHERE id = $2`, tableName)

	_, err := r.DB.Exec(query, domain.Denied, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgReportRepository) Delete(id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, tableName)

	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgReportRepository) checkIfRecordExists(id, table string) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)`, table)

	var exists bool
	err := r.DB.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *PgReportRepository) checkIfReportMatchesStatus(id string, status domain.ReportStatus) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1 AND status = $2)`, tableName)

	var matches bool
	err := r.DB.QueryRow(query, id, status).Scan(&matches)
	if err != nil {
		return false, err
	}

	return matches, nil
}

func (r *PgReportRepository) getFullReportById(id string) (*domain.Report, error) {
	return r.getFullReport(id, nil)
}

func (r *PgReportRepository) getFullReportByIdAndStatus(id string, status domain.ReportStatus) (*domain.Report, error) {
	return r.getFullReport(id, &status)
}

func (r *PgReportRepository) getFullReport(id string, status *domain.ReportStatus) (*domain.Report, error) {
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
		row = r.DB.QueryRow(baseQuery, id, *status)
	} else {
		row = r.DB.QueryRow(baseQuery, id)
	}

	return r.ScanReportRow(row)
}

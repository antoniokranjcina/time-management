package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time-management/internal/shared/util"
	"time-management/internal/user"
	"time-management/internal/user/employee/domain"
)

type PgEmployeeRepository struct {
	DB *sql.DB
}

func NewPgEmployeeRepository(db *sql.DB) *PgEmployeeRepository {
	repository := &PgEmployeeRepository{DB: db}

	err := repository.createEmployeeTable()
	if err != nil {
		panic(err)
		return nil
	}

	return repository
}

func (r *PgEmployeeRepository) createEmployeeTable() error {
	query := fmt.Sprintf(`create table if not exists %s (
		id varchar(50) primary key,
		first_name varchar(50),
		last_name varchar(50),
		email varchar(50),
    	role varchar(50),
		password_hashed varchar,
		created_at serial,
		active boolean
	)`, user.TableName)

	_, err := r.DB.Exec(query)
	return err
}

func (r *PgEmployeeRepository) Save(employee *user.User) (*user.User, error) {
	if exists, err := r.isEmailTaken(employee.Email); err != nil {
		return nil, err
	} else if exists {
		return nil, domain.ErrEmailTaken
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (id, first_name, last_name, email, role, password_hashed, created_at, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id, first_name, last_name, email, role, password_hashed, created_at, active
	`, user.TableName)

	row := r.DB.QueryRow(
		query,
		employee.Id,
		employee.FirstName,
		employee.LastName,
		employee.Email,
		employee.Role,
		employee.PasswordHash,
		employee.CreatedAt,
		employee.Active,
	)

	savedEmployee, err := user.ScanUserRow(row)
	if err != nil {
		return nil, err
	}

	return savedEmployee, nil
}

func (r *PgEmployeeRepository) GetAll() ([]user.User, error) {
	query := fmt.Sprintf(`SELECT * FROM %s`, user.TableName)

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employees, err := user.ScanUserRows(rows)
	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (r *PgEmployeeRepository) GetById(id string) (*user.User, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, user.TableName)

	row := r.DB.QueryRow(query, id)
	employee, err := user.ScanUserRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, util.NewValidationError(domain.ErrEmployeeNotFound)
		}
		return nil, err
	}

	return employee, nil
}

func (r *PgEmployeeRepository) Update(id, firstName, lastName string) (*user.User, error) {
	query := fmt.Sprintf(`
		UPDATE %s SET first_name = $1, last_name = $2 
	 	WHERE id = $3 
		RETURNING id, first_name, last_name, email, role, password_hashed, created_at, active
	`, user.TableName)

	row := r.DB.QueryRow(query, firstName, lastName, id)
	employee, err := user.ScanUserRow(row)
	if err != nil {
		return nil, err
	}

	return employee, nil
}

func (r *PgEmployeeRepository) ChangePassword(id, password string) error {
	query := fmt.Sprintf(`UPDATE %s SET password_hashed = $1 WHERE id = $2`, user.TableName)

	_, err := r.DB.Exec(query, password, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgEmployeeRepository) ChangeEmail(id, email string) error {
	currentMail, err := r.getEmailById(id)
	if err != nil {
		return err
	}
	if currentMail == email {
		return util.NewValidationError(domain.ErrEmailTaken)
	}

	query := fmt.Sprintf(`UPDATE %s SET email = $1 WHERE id = $2`, user.TableName)

	_, err = r.DB.Exec(query, email, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgEmployeeRepository) ToggleStatus(id string, status bool) (bool, error) {
	query := fmt.Sprintf(`UPDATE %s SET active = $1 WHERE id = $2 RETURNING active`, user.TableName)

	var newStatus bool
	err := r.DB.QueryRow(query, status, id).Scan(&newStatus)
	if err != nil {
		return false, err
	}

	return newStatus, nil
}

func (r *PgEmployeeRepository) Delete(id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, user.TableName)

	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgEmployeeRepository) isEmailTaken(email string) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE email = $1)`, user.TableName)

	var exists bool
	err := r.DB.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *PgEmployeeRepository) getEmailById(id string) (string, error) {
	query := fmt.Sprintf(`SELECT email FROM %s WHERE id = $1`, user.TableName)

	email := ""
	row := r.DB.QueryRow(query, id)
	err := row.Scan(&email)
	if err != nil {
		return "", err
	}

	return email, nil
}

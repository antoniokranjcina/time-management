package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time-management/internal/employee/domain"
	"time-management/internal/shared/util"
)

const TableName = "employees"

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
		password varchar(50),
		created_at serial,
		active boolean
	)`, TableName)

	_, err := r.DB.Exec(query)
	return err
}

func (r *PgEmployeeRepository) Save(employee *domain.Employee) (*domain.Employee, error) {
	if exists, err := r.isEmailTaken(employee.Email); err != nil {
		return nil, err
	} else if exists {
		return nil, domain.ErrEmailTaken
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (id, first_name, last_name, email, password, created_at, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id, first_name, last_name, email, password, created_at, active
	`, TableName)

	row := r.DB.QueryRow(
		query,
		employee.Id,
		employee.FirstName,
		employee.LastName,
		employee.Email,
		employee.Password,
		employee.CreatedAt,
		employee.Active,
	)

	savedEmployee, err := scanEmployeeRow(row)
	if err != nil {
		return nil, err
	}

	return savedEmployee, nil
}

func (r *PgEmployeeRepository) GetAll() ([]domain.Employee, error) {
	query := fmt.Sprintf(`SELECT * FROM %s`, TableName)

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employees, err := scanEmployeeRows(rows)
	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (r *PgEmployeeRepository) GetById(id string) (*domain.Employee, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, TableName)

	row := r.DB.QueryRow(query, id)
	employee, err := scanEmployeeRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, util.NewValidationError(domain.ErrEmployeeNotFound)
		}
		return nil, err
	}

	return employee, nil
}

func (r *PgEmployeeRepository) Update(id, firstName, lastName string) (*domain.Employee, error) {
	query := fmt.Sprintf(`
		UPDATE %s SET first_name = $1, last_name = $2 
	 	WHERE id = $3 
		RETURNING id, first_name, last_name, email, password, created_at, active
	`, TableName)

	row := r.DB.QueryRow(query, firstName, lastName, id)
	employee, err := scanEmployeeRow(row)
	if err != nil {
		return nil, err
	}

	return employee, nil
}

func (r *PgEmployeeRepository) ChangePassword(id, password string) error {
	query := fmt.Sprintf(`UPDATE %s SET password = $1 WHERE id = $2`, TableName)

	_, err := r.DB.Exec(query, password, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgEmployeeRepository) ChangeEmail(id, email string) error {
	query := fmt.Sprintf(`UPDATE %s SET email = $1 WHERE id = $2`, TableName)

	_, err := r.DB.Exec(query, email, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgEmployeeRepository) ToggleStatus(id string, status bool) (bool, error) {
	query := fmt.Sprintf(`UPDATE %s SET active = $1 WHERE id = $2 RETURNING active`, TableName)

	var newStatus bool
	err := r.DB.QueryRow(query, status, id).Scan(&newStatus)
	if err != nil {
		return false, err
	}

	return newStatus, nil
}

func (r *PgEmployeeRepository) Delete(id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, TableName)

	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgEmployeeRepository) isEmailTaken(email string) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE email = $1)`, TableName)

	var exists bool
	err := r.DB.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func scanEmployeeRow(row *sql.Row) (*domain.Employee, error) {
	employee := &domain.Employee{}
	err := row.Scan(
		&employee.Id,
		&employee.FirstName,
		&employee.LastName,
		&employee.Email,
		&employee.Password,
		&employee.CreatedAt,
		&employee.Active,
	)
	if err != nil {
		return nil, err
	}

	return employee, nil
}

func scanEmployeeRows(rows *sql.Rows) ([]domain.Employee, error) {
	var employees []domain.Employee

	for rows.Next() {
		var employee domain.Employee
		err := rows.Scan(
			&employee.Id,
			&employee.FirstName,
			&employee.LastName,
			&employee.Email,
			&employee.Password,
			&employee.CreatedAt,
			&employee.Active,
		)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}

	return employees, nil
}

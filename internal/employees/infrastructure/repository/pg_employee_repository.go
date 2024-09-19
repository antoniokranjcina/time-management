package repository

import (
	"database/sql"
	"errors"
	"time-management/internal/employees/domain"
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
	query := `create table if not exists employees (
		id varchar(50) primary key,
		first_name varchar(50),
		last_name varchar(50),
		email varchar(50),
		password varchar(50),
		created_at serial,
		active boolean
	)`

	_, err := r.DB.Exec(query)
	return err
}

func (r *PgEmployeeRepository) Save(employee *domain.Employee) (*domain.Employee, error) {
	var exists bool
	emailCheckQuery := `SELECT EXISTS(SELECT 1 FROM employees WHERE email = $1)`
	err := r.DB.QueryRow(emailCheckQuery, employee.Email).Scan(&exists)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, domain.ErrEmailTaken
	}

	query := `
		INSERT INTO employees (id, first_name, last_name, email, password, created_at, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id, first_name, last_name, email, password, created_at, active
	`

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
	query := `SELECT * FROM employees`
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
	query := `SELECT * FROM employees WHERE id = $1`
	row := r.DB.QueryRow(query, id)

	employee, err := scanEmployeeRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrEmployeeNotFound
		}
		return nil, err
	}

	return employee, nil
}

func (r *PgEmployeeRepository) Update(id, firstName, lastName string) (*domain.Employee, error) {
	query := `
		UPDATE employees SET first_name = $1, last_name = $2 
	 	WHERE id = $3 
		RETURNING id, first_name, last_name, email, password, created_at, active
	`

	row := r.DB.QueryRow(query, firstName, lastName, id)
	employee, err := scanEmployeeRow(row)
	if err != nil {
		return nil, err
	}

	return employee, nil
}

func (r *PgEmployeeRepository) ChangePassword(id, password string) error {
	query := `UPDATE employees SET password = $1 WHERE id = $2`

	_, err := r.DB.Exec(query, password, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgEmployeeRepository) ChangeEmail(id, email string) error {
	query := `UPDATE employees SET email = $1 WHERE id = $2`

	_, err := r.DB.Exec(query, email, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgEmployeeRepository) ToggleStatus(id string, status bool) (bool, error) {
	query := `UPDATE employees SET active = $1 WHERE id = $2 RETURNING active`

	var newStatus bool
	err := r.DB.QueryRow(query, status, id).Scan(&newStatus)
	if err != nil {
		return false, err
	}

	return newStatus, nil
}

func (r *PgEmployeeRepository) Delete(id string) error {
	query := `DELETE FROM employees WHERE id = $1`

	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
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

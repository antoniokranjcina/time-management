package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time-management/internal/shared/util"
	"time-management/internal/user/domain"
)

const TableName = "users"

type PgUserRepository struct {
	DB *sql.DB
}

func NewPgUsersRepository(db *sql.DB) *PgUserRepository {
	repository := &PgUserRepository{DB: db}

	err := repository.createUsersTable()
	if err != nil {
		panic(err)
		return nil
	}

	return repository
}

func (r *PgUserRepository) createUsersTable() error {
	query := fmt.Sprintf(`create table if not exists %s (
		id varchar(50) primary key,
		first_name varchar(50),
		last_name varchar(50),
		email varchar(50),
    	role varchar(50),
		password_hashed varchar,
		created_at serial,
		active boolean
	)`, TableName)

	_, err := r.DB.Exec(query)
	return err
}

func (r *PgUserRepository) Save(user *domain.User) (*domain.User, error) {
	if exists, err := r.isEmailTaken(user.Email); err != nil {
		return nil, err
	} else if exists {
		return nil, domain.ErrEmailTaken
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (id, first_name, last_name, email, role, password_hashed, created_at, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id, first_name, last_name, email, role, password_hashed, created_at, active
	`, TableName)

	row := r.DB.QueryRow(
		query,
		user.Id,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Role,
		user.PasswordHash,
		user.CreatedAt,
		user.Active,
	)

	savedEmployee, err := ScanUserRow(row)
	if err != nil {
		return nil, err
	}

	return savedEmployee, nil
}

func (r *PgUserRepository) GetAllWithRole(role string) ([]domain.User, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE role = $1`, TableName)

	rows, err := r.DB.Query(query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employees, err := ScanUserRows(rows)
	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (r *PgUserRepository) GetByIdWithRole(id, role string) (*domain.User, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1 AND role = $2`, TableName)

	row := r.DB.QueryRow(query, id, role)
	user, err := ScanUserRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, util.NewValidationError(domain.ErrUserNotFound)
		}
		return nil, err
	}

	return user, nil
}

func (r *PgUserRepository) Update(id, firstName, lastName string) (*domain.User, error) {
	query := fmt.Sprintf(`
		UPDATE %s SET first_name = $1, last_name = $2 
	 	WHERE id = $3 
		RETURNING id, first_name, last_name, email, role, password_hashed, created_at, active
	`, TableName)

	row := r.DB.QueryRow(query, firstName, lastName, id)
	user, err := ScanUserRow(row)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *PgUserRepository) ChangePassword(id, password string) error {
	query := fmt.Sprintf(`UPDATE %s SET password_hashed = $1 WHERE id = $2`, TableName)

	_, err := r.DB.Exec(query, password, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgUserRepository) ChangeEmail(id, email string) error {
	currentMail, err := r.getEmailById(id)
	if err != nil {
		return err
	}
	if currentMail == email {
		return util.NewValidationError(domain.ErrEmailTaken)
	}

	query := fmt.Sprintf(`UPDATE %s SET email = $1 WHERE id = $2`, TableName)

	_, err = r.DB.Exec(query, email, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgUserRepository) ToggleStatus(id string, status bool) (bool, error) {
	query := fmt.Sprintf(`UPDATE %s SET active = $1 WHERE id = $2 RETURNING active`, TableName)

	var newStatus bool
	err := r.DB.QueryRow(query, status, id).Scan(&newStatus)
	if err != nil {
		return false, err
	}

	return newStatus, nil
}

func (r *PgUserRepository) Delete(id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, TableName)

	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgUserRepository) isEmailTaken(email string) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE email = $1)`, TableName)

	var exists bool
	err := r.DB.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *PgUserRepository) getEmailById(id string) (string, error) {
	query := fmt.Sprintf(`SELECT email FROM %s WHERE id = $1`, TableName)

	email := ""
	row := r.DB.QueryRow(query, id)
	err := row.Scan(&email)
	if err != nil {
		return "", err
	}

	return email, nil
}

package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
	"time-management/internal/shared/util"
	"time-management/internal/user/domain"
	"time-management/internal/user/role"
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

	err = r.createSuperAdmin()

	return err
}

func (r *PgUserRepository) createSuperAdmin() error {
	// First, check if a super-admin already exists
	checkQuery := fmt.Sprintf(`SELECT id FROM %s WHERE role = $1 LIMIT 1`, TableName)

	var superAdminId string
	err := r.DB.QueryRow(checkQuery, role.SuperAdmin.String()).Scan(&superAdminId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) { // If there's an error that's not "no rows"
		return fmt.Errorf("failed to check for existing super admin: %v", err)
	}

	// If super-admin already exists, return without doing anything
	if superAdminId != "" {
		fmt.Println("Super admin already exists. Skipping creation.")
		return nil
	}

	// Otherwise, proceed with creating a new super-admin
	query := fmt.Sprintf(`
		INSERT INTO %s (id, first_name, last_name, email, role, password_hashed, created_at, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
	`, TableName)

	// Fetch email and password from environment variables
	password := os.Getenv("SUPER_ADMIN_PASSWORD")
	email := os.Getenv("SUPER_ADMIN_EMAIL")

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return domain.ErrFailedToHashPassword
	}

	// Insert the new super admin user
	_, err = r.DB.Exec(
		query,
		uuid.New().String(),
		"Super",
		"Admin",
		email,
		role.SuperAdmin.String(),
		string(hash),
		uint64(time.Now().Unix()),
		true,
	)
	if err != nil {
		return fmt.Errorf("failed to create super admin: %v", err)
	}

	fmt.Println("Super admin created successfully")
	return nil
}

func (r *PgUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	if exists, err := r.isEmailTaken(user.Email); err != nil {
		return nil, err
	} else if exists {
		return nil, util.NewValidationError(domain.ErrEmailTaken)
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (id, first_name, last_name, email, role, password_hashed, created_at, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id, first_name, last_name, email, role, password_hashed, created_at, active
	`, TableName)

	row := r.DB.QueryRowContext(
		ctx,
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

	savedUser, err := ScanUserRow(row)
	if err != nil {
		return nil, err
	}

	return savedUser, nil
}

func (r *PgUserRepository) GetAllWithRole(ctx context.Context, role string) ([]domain.User, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE role = $1`, TableName)

	rows, err := r.DB.QueryContext(ctx, query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users, err := ScanUserRows(rows)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *PgUserRepository) GetByIdWithRole(ctx context.Context, id, role string) (*domain.User, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1 AND role = $2`, TableName)

	row := r.DB.QueryRowContext(ctx, query, id, role)
	user, err := ScanUserRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, util.NewValidationError(domain.ErrUserNotFound)
		}
		return nil, err
	}

	return user, nil
}

func (r *PgUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE email = $1`, TableName)

	row := r.DB.QueryRowContext(ctx, query, email)
	user, err := ScanUserRow(row)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *PgUserRepository) Update(ctx context.Context, id, firstName, lastName string) (*domain.User, error) {
	query := fmt.Sprintf(`
		UPDATE %s SET first_name = $1, last_name = $2 
	 	WHERE id = $3 
		RETURNING id, first_name, last_name, email, role, password_hashed, created_at, active
	`, TableName)

	row := r.DB.QueryRowContext(ctx, query, firstName, lastName, id)
	user, err := ScanUserRow(row)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *PgUserRepository) ChangePassword(ctx context.Context, id, password string) error {
	query := fmt.Sprintf(`UPDATE %s SET password_hashed = $1 WHERE id = $2`, TableName)

	_, err := r.DB.ExecContext(ctx, query, password, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgUserRepository) ChangeEmail(ctx context.Context, id, email string) error {
	currentMail, err := r.getEmailById(id)
	if err != nil {
		return err
	}
	if currentMail == email {
		return util.NewValidationError(domain.ErrEmailTaken)
	}

	query := fmt.Sprintf(`UPDATE %s SET email = $1 WHERE id = $2`, TableName)

	_, err = r.DB.ExecContext(ctx, query, email, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgUserRepository) ToggleStatus(ctx context.Context, id string, status bool) (bool, error) {
	query := fmt.Sprintf(`UPDATE %s SET active = $1 WHERE id = $2 RETURNING active`, TableName)

	var newStatus bool
	err := r.DB.QueryRowContext(ctx, query, status, id).Scan(&newStatus)
	if err != nil {
		return false, err
	}

	return newStatus, nil
}

func (r *PgUserRepository) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, TableName)

	_, err := r.DB.ExecContext(ctx, query, id)
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

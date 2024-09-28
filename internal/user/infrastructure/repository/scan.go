package repository

import (
	"database/sql"
	"time-management/internal/user/domain"
)

func ScanUserRow(row *sql.Row) (*domain.User, error) {
	user := &domain.User{}
	err := row.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Role,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.Active,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func ScanUserRows(rows *sql.Rows) ([]domain.User, error) {
	var users []domain.User

	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Role,
			&user.PasswordHash,
			&user.CreatedAt,
			&user.Active,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

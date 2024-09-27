package user

import (
	"database/sql"
)

func ScanUserRow(row *sql.Row) (*User, error) {
	user := &User{}
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

func ScanUserRows(rows *sql.Rows) ([]User, error) {
	var users []User

	for rows.Next() {
		var user User
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

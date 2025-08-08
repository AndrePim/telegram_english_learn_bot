package repository

import (
	"database/sql"
	"fmt"
)

// User представляет собой структуру пользователя
type UserRepository struct {
	db *sql.DB
}

// User представляет собой структуру пользователя
func NewUserRepository(database *Database) *UserRepository {
	return &UserRepository{db: database.db}
}

// CreateOrUpdateUser создает или обновляет пользователя
func (r *UserRepository) CreateOrUpdateUser(user *User) error {
	query := `
		INSERT INTO users (id, username, first_name, last_name, state)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			username = EXCLUDED.username,
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name
	`

	_, err := r.db.Exec(query, user.ID, user.Username, user.FirstName, user.LastName, user.State)
	if err != nil {
		return fmt.Errorf("failed to create or update user: %w", err)
	}

	return nil
}

// GetUser получает пользователя по ID
func (r *UserRepository) GetUser(userID int64) (*User, error) {
	query := `
		SELECT id, username, first_name, last_name, state, created_at
		FROM users WHERE id = $1
	`

	user := &User{}
	err := r.db.QueryRow(query, userID).Scan(
		&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.State, &user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Пользователь не найден
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// UpdateUserState обновляет состояние пользователя
func (r *UserRepository) UpdateUserState(userID int64, state string) error {
	query := `UPDATE users SET state = $1 WHERE id = $2`

	_, err := r.db.Exec(query, state, userID)
	if err != nil {
		return fmt.Errorf("failed to update user state: %w", err)
	}

	return nil
}

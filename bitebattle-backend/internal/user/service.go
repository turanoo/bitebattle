package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	DB *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) CreateUser(ctx context.Context, u *User) (*User, error) {
	u.ID = uuid.NewString()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO users (id, email, name, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, u.ID, u.Email, u.Name, u.PasswordHash, u.CreatedAt, u.UpdatedAt)

	return u, err
}

func (s *Service) GetUserByID(ctx context.Context, id string) (*User, error) {
	var u User
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, email, name, password_hash, created_at, updated_at
		FROM users WHERE id = $1
	`, id)

	if err := row.Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found with this ID
		}
		return nil, err
	}

	return &u, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, email, name, password_hash, created_at, updated_at
		FROM users WHERE email = $1
	`, email)

	if err := row.Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found with this email
		}
		return nil, err
	}

	return &u, nil
}

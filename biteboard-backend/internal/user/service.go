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

func (s *Service) CreateUser(ctx context.Context, u *User) error {
	u.ID = uuid.NewString()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO users (id, email, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, u.ID, u.Email, u.Name, u.CreatedAt, u.UpdatedAt)

	return err
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, email, name, created_at, updated_at
		FROM users WHERE email = $1
	`, email)

	var u User
	if err := row.Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}

	return &u, nil
}

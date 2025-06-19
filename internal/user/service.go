package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/turanoo/bitebattle/pkg/db"
)

var ErrUserExists = errors.New("user with this email already exists")

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
		INSERT INTO users (id, email, name, password_hash, phone_number, profile_pic_url, bio, last_login_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, u.ID, u.Email, u.Name, u.PasswordHash, u.PhoneNumber, u.ProfilePicURL, u.Bio, u.LastLoginAt, u.CreatedAt, u.UpdatedAt)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, ErrUserExists
		}
		return nil, err
	}

	return u, nil
}

func (s *Service) GetUserByID(ctx context.Context, id string) (*User, error) {
	var u User
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, email, name, password_hash, phone_number, profile_pic_url, bio, last_login_at, created_at, updated_at
		FROM users WHERE id = $1
	`, id)

	err := db.ScanOne(row, &u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.PhoneNumber, &u.ProfilePicURL, &u.Bio, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if u.ID == "" {
		return nil, sql.ErrNoRows
	}
	return &u, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, email, name, password_hash, phone_number, profile_pic_url, bio, last_login_at, created_at, updated_at
		FROM users WHERE email = $1
	`, email)

	err := db.ScanOne(row, &u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.PhoneNumber, &u.ProfilePicURL, &u.Bio, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if u.ID == "" {
		return nil, sql.ErrNoRows
	}
	return &u, nil
}

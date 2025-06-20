package account

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/turanoo/bitebattle/internal/user"
	"github.com/turanoo/bitebattle/pkg/db"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword = errors.New("current password is incorrect")
	ErrEmailExists     = errors.New("user with this email already exists")
)

type Service struct {
	DB *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{DB: db}
}

type UserProfile struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	PhoneNumber   *string   `json:"phone_number,omitempty"`
	ProfilePicURL *string   `json:"profile_pic_url,omitempty"`
	Bio           *string   `json:"bio,omitempty"`
	LastLoginAt   *string   `json:"last_login_at,omitempty"`
}

func (s *Service) GetUserProfile(userID uuid.UUID) (*UserProfile, error) {
	row := s.DB.QueryRow(`SELECT id, name, email, phone_number, profile_pic_url, bio, last_login_at FROM users WHERE id = $1`, userID)

	var profile UserProfile
	err := db.ScanOne(row, &profile.ID, &profile.Name, &profile.Email, &profile.PhoneNumber, &profile.ProfilePicURL, &profile.Bio, &profile.LastLoginAt)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func isEmptyUpdateFields(name, email, password *string) bool {
	return name == nil && email == nil && password == nil
}

func (s *Service) UpdateUserProfile(userID uuid.UUID, name, email, currentPassword, newPassword *string) error {
	if isEmptyUpdateFields(name, email, newPassword) {
		return errors.New("no fields provided for updating user profile")
	}

	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if name != nil {
		setClauses = append(setClauses, "name = $"+strconv.Itoa(argIdx))
		args = append(args, *name)
		argIdx++
	}
	if email != nil {
		setClauses = append(setClauses, "email = $"+strconv.Itoa(argIdx))
		args = append(args, *email)
		argIdx++
	}

	if currentPassword != nil && newPassword != nil {

		var storedHash string
		err := s.DB.QueryRow(`SELECT password_hash FROM users WHERE id = $1`, userID).Scan(&storedHash)
		if err != nil {
			if err == sql.ErrNoRows {
				return sql.ErrNoRows
			}
			return err
		}

		if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(*currentPassword)); err != nil {
			return ErrInvalidPassword
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(*newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		setClauses = append(setClauses, "password_hash = $"+strconv.Itoa(argIdx))
		args = append(args, string(hashed))
		argIdx++
	} else if currentPassword != nil || newPassword != nil {
		return errors.New("both currentPassword and newPassword must be provided to update password")
	}

	if len(setClauses) == 0 {
		return errors.New("no valid fields provided for update")
	}

	query := "UPDATE users SET " + strings.Join(setClauses, ", ") + " WHERE id = $" + strconv.Itoa(argIdx)
	args = append(args, userID)

	result, err := s.DB.Exec(query, args...)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID uuid.UUID, name, email string) error {
	_, err := s.DB.ExecContext(ctx, `
		UPDATE users SET name = $1, email = $2, updated_at = NOW()
		WHERE id = $3
	`, name, email, userID)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrEmailExists
		}
		return err
	}
	return nil
}

func (s *Service) GetProfile(ctx context.Context, userID string) (*user.User, error) {
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE id = $1
	`, userID)

	var u user.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

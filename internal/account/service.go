package account

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/turanoo/bitebattle/pkg/db"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword = errors.New("current password is incorrect")
)

type Service struct {
	DB *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{DB: db}
}

type UserProfile struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

func (s *Service) GetUserProfile(userID uuid.UUID) (*UserProfile, error) {
	row := s.DB.QueryRow(`SELECT id, name, email FROM users WHERE id = $1`, userID)

	var profile UserProfile
	err := db.ScanOne(row, &profile.ID, &profile.Name, &profile.Email)
	if err != nil {
		return nil, err
	}
	if profile.ID == uuid.Nil {
		return nil, sql.ErrNoRows
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

	// Handle password update
	if currentPassword != nil && newPassword != nil {
		// Fetch current password hash
		var storedHash string
		err := s.DB.QueryRow(`SELECT password_hash FROM users WHERE id = $1`, userID).Scan(&storedHash)
		if err != nil {
			if err == sql.ErrNoRows {
				return sql.ErrNoRows
			}
			return err
		}
		// Validate current password
		if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(*currentPassword)); err != nil {
			return ErrInvalidPassword
		}
		// Hash new password
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

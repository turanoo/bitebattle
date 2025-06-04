package account

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/turanoo/bitebattle/bitebattle-backend/pkg/db"
	"golang.org/x/crypto/bcrypt"
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

func (s *Service) UpdateUserProfile(userID uuid.UUID, name, email, password *string) error {
	if isEmptyUpdateFields(name, email, password) {
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
	if password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		setClauses = append(setClauses, "password_hash = $"+strconv.Itoa(argIdx))
		args = append(args, string(hashed))
		argIdx++
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

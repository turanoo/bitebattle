package account

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/google/uuid"
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
	if err := row.Scan(&profile.ID, &profile.Name, &profile.Email); err != nil {
		return nil, err
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

	_, err := s.DB.Exec(query, args...)
	return err
}

type GroupSummary struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Role string    `json:"role"` // "owner" or "member"
}

func (s *Service) GetUserGroups(userID uuid.UUID) ([]GroupSummary, error) {
	rows, err := s.DB.Query(`
		SELECT g.id, g.name, 
			CASE 
				WHEN g.created_by = $1 THEN 'owner'
				ELSE 'member'
			END as role
		FROM groups g
		LEFT JOIN group_members gm ON g.id = gm.group_id
		WHERE g.created_by = $1 OR gm.user_id = $1
		GROUP BY g.id
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []GroupSummary
	for rows.Next() {
		var g GroupSummary
		if err := rows.Scan(&g.ID, &g.Name, &g.Role); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}

	return groups, nil
}

type PollSummary struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"restaurant_name"`
}

func (s *Service) GetUserPolls(userID uuid.UUID) ([]PollSummary, error) {
	rows, err := s.DB.Query(`
		SELECT DISTINCT p.id, po.name
		FROM polls p
		JOIN poll_votes v ON p.id = v.poll_id
		JOIN poll_options po ON v.option_id = po.id
		WHERE v.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var polls []PollSummary
	for rows.Next() {
		var p PollSummary
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			return nil, err
		}
		polls = append(polls, p)
	}

	return polls, nil
}

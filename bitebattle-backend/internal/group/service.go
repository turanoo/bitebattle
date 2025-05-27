package group

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/turanoo/bitebattle/bitebattle-backend/pkg/utils"
)

type Service struct {
	DB *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) CreateGroup(name string, createdBy uuid.UUID) (*Group, error) {
	id := uuid.New()
	inviteCode := utils.GenerateRandomString(8) // Make sure you have this util
	now := time.Now()

	_, err := s.DB.Exec(`
		INSERT INTO groups (id, name, created_by, invite_code, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, id, name, createdBy, inviteCode, now, now)
	if err != nil {
		return nil, err
	}

	return &Group{
		ID:         id,
		Name:       name,
		CreatedBy:  createdBy,
		InviteCode: inviteCode,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func (s *Service) GetGroupByID(groupID uuid.UUID) (*Group, error) {
	row := s.DB.QueryRow(`
		SELECT id, name, created_by, invite_code, created_at, updated_at
		FROM groups WHERE id = $1
	`, groupID)

	var group Group
	err := row.Scan(
		&group.ID, &group.Name, &group.CreatedBy, &group.InviteCode, &group.CreatedAt, &group.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &group, nil
}

func (s *Service) JoinGroupByInviteCode(code, userId string) (*Group, error) {
	// Step 1: Fetch the group by invite code
	row := s.DB.QueryRow(`
		SELECT id, name, created_by, invite_code, created_at, updated_at
		FROM groups WHERE invite_code = $1
	`, code)

	var group Group
	err := row.Scan(
		&group.ID, &group.Name, &group.CreatedBy, &group.InviteCode, &group.CreatedAt, &group.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid invite code")
		}
		return nil, err
	}

	// Step 2: Insert into group_members table
	_, err = s.DB.Exec(`
		INSERT INTO group_members (group_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (group_id, user_id) DO NOTHING
	`, group.ID, userId)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

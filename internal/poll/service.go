package poll

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/turanoo/bitebattle/pkg/db"
	"github.com/turanoo/bitebattle/pkg/utils"
)

type Service struct {
	DB *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) CreatePoll(name string, createdBy uuid.UUID) (*Poll, error) {
	id := uuid.New()
	now := time.Now()
	inviteCode := utils.GenerateRandomString(8)

	_, err := s.DB.Exec(`
		INSERT INTO polls (id, name, invite_code, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, id, name, inviteCode, createdBy, now, now)

	if err != nil {
		return nil, err
	}

	poll := Poll{
		ID:         id,
		Name:       name,
		InviteCode: inviteCode,
		Role:       "owner", // Creator is always the owner
		CreatedBy:  createdBy,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	_, err = s.DB.Exec(`
		INSERT INTO polls_members (poll_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (poll_id, user_id) DO NOTHING
	`, poll.ID, createdBy)
	if err != nil {
		return nil, err
	}

	return &poll, nil
}

func (s *Service) GetPolls(userID uuid.UUID) ([]Poll, error) {
	rows, err := s.DB.Query(`
		SELECT p.id, p.name, p.invite_code, p.created_by, p.created_at, p.updated_at,
			CASE 
				WHEN p.created_by = $1 THEN 'owner'
				WHEN EXISTS (
					SELECT 1 FROM polls_members WHERE poll_id = p.id AND user_id = $1
				) THEN 'member'
				ELSE NULL
			END as role
		FROM polls p
		LEFT JOIN polls_members pm ON p.id = pm.poll_id
		WHERE p.created_by = $1 OR pm.user_id = $1
		GROUP BY p.id
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	polls := []Poll{}
	for rows.Next() {
		var poll Poll
		if err := rows.Scan(
			&poll.ID,
			&poll.Name,
			&poll.InviteCode,
			&poll.CreatedBy,
			&poll.CreatedAt,
			&poll.UpdatedAt,
			&poll.Role,
		); err != nil {
			return nil, err
		}

		memberRows, err := s.DB.Query(`
			SELECT user_id FROM polls_members WHERE poll_id = $1
		`, poll.ID)
		if err != nil {
			return nil, err
		}
		members := []uuid.UUID{}
		for memberRows.Next() {
			var memberID uuid.UUID
			if err := memberRows.Scan(&memberID); err != nil {
				memberRows.Close()
				return nil, err
			}
			members = append(members, memberID)
		}
		memberRows.Close()
		poll.Members = members

		polls = append(polls, poll)
	}

	return polls, nil
}

func (s *Service) GetPoll(pollID, userId uuid.UUID) (*Poll, error) {
	row := s.DB.QueryRow(`
		SELECT id, name, invite_code, created_by, created_at, updated_at,
			CASE 
				WHEN created_by = $2 THEN 'owner'
				WHEN EXISTS (
					SELECT 1 FROM polls_members WHERE poll_id = polls.id AND user_id = $2
				) THEN 'member'
				ELSE NULL
			END as role
		FROM polls
		WHERE id = $1
	`, pollID, userId)

	var poll Poll
	err := db.ScanOne(row, &poll.ID, &poll.Name, &poll.InviteCode, &poll.CreatedBy, &poll.CreatedAt, &poll.UpdatedAt, &poll.Role)
	if err != nil {
		return nil, err
	}
	if poll.ID == uuid.Nil || poll.Role == "" {
		return nil, sql.ErrNoRows
	}

	memberRows, err := s.DB.Query(`
		SELECT user_id FROM polls_members WHERE poll_id = $1
	`, pollID)
	if err != nil {
		return nil, err
	}
	defer memberRows.Close()

	members := []uuid.UUID{}
	for memberRows.Next() {
		var memberID uuid.UUID
		if err := memberRows.Scan(&memberID); err != nil {
			return nil, err
		}
		members = append(members, memberID)
	}
	poll.Members = members

	return &poll, nil
}

func (s *Service) DeletePoll(pollID uuid.UUID) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	_, err = tx.Exec(`
		DELETE FROM poll_votes WHERE option_id IN (
			SELECT id FROM poll_options WHERE poll_id = $1
		)
	`, pollID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM poll_options WHERE poll_id = $1
	`, pollID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM polls_members WHERE poll_id = $1
	`, pollID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM polls WHERE id = $1
	`, pollID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *Service) UpdatePoll(pollID uuid.UUID, name string) (*Poll, error) {
	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if name != "" {
		setClauses = append(setClauses, "name = $"+strconv.Itoa(argIdx))
		args = append(args, name)
		argIdx++
	}

	if len(setClauses) == 0 {
		return nil, errors.New("no fields provided for updating poll")
	}

	query := "UPDATE polls SET " + strings.Join(setClauses, ", ") + ", updated_at = NOW() WHERE id = $" + strconv.Itoa(argIdx)
	args = append(args, pollID)

	_, err := s.DB.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return s.GetPoll(pollID, uuid.Nil) // Return updated poll without user context
}

func (s *Service) JoinPoll(inviteCode string, userId uuid.UUID) (*Poll, error) {
	row := s.DB.QueryRow(`
		SELECT id, name, created_by, invite_code, created_at, updated_at
		FROM polls WHERE invite_code = $1
	`, inviteCode)

	var poll Poll
	err := row.Scan(
		&poll.ID, &poll.Name, &poll.CreatedBy, &poll.InviteCode, &poll.CreatedAt, &poll.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid invite code")
		}
		return nil, err
	}

	_, err = s.DB.Exec(`
		INSERT INTO polls_members (poll_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (poll_id, user_id) DO NOTHING
	`, poll.ID, userId)
	if err != nil {
		return nil, err
	}

	poll.Role = "member" // Default role for joined users

	return &poll, nil
}

func (s *Service) AddOption(pollID uuid.UUID, restaurantID, name, imageURL, menuURL string) (*PollOption, error) {
	id := uuid.New()

	_, err := s.DB.Exec(`
        INSERT INTO poll_options (id, poll_id, restaurant_id, name, image_url, menu_url)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, id, pollID, restaurantID, name, imageURL, menuURL)

	if err != nil {
		return nil, err
	}

	return &PollOption{
		ID:           id,
		PollID:       pollID,
		RestaurantID: restaurantID,
		Name:         name,
		ImageURL:     imageURL,
		MenuURL:      menuURL,
	}, nil
}

func (s *Service) CastVote(pollID, optionID, userID uuid.UUID) (*PollVote, error) {
	id := uuid.New()

	_, err := s.DB.Exec(`
		INSERT INTO poll_votes (id, poll_id, option_id, user_id)
		VALUES ($1, $2, $3, $4)
	`, id, pollID, optionID, userID)

	if err != nil {
		return nil, err
	}

	return &PollVote{
		ID:       id,
		PollID:   pollID,
		OptionID: optionID,
		UserID:   userID,
	}, nil
}

func (s *Service) RemoveVote(pollID, optionID, userID uuid.UUID) error {
	result, err := s.DB.Exec(`
		DELETE FROM poll_votes WHERE poll_id = $1 AND option_id = $2 AND user_id = $3
	`, pollID, optionID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("vote not found")
	}

	return nil
}

func (s *Service) GetResults(pollID uuid.UUID) ([]PollResult, error) {
	rows, err := s.DB.Query(`
		SELECT o.id, o.name, COUNT(v.id) as votes
		FROM poll_options o
		LEFT JOIN poll_votes v ON o.id = v.option_id
		WHERE o.poll_id = $1
		GROUP BY o.id
		ORDER BY votes DESC
	`, pollID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []PollResult{}
	for rows.Next() {
		var res PollResult
		var optionID uuid.UUID
		var optionName string
		var voteCount int

		voterIds := []uuid.UUID{}

		if err := rows.Scan(&optionID, &optionName, &voteCount); err != nil {
			return nil, err
		}

		voterRows, err := s.DB.Query(`
			SELECT user_id FROM poll_votes WHERE option_id = $1
		`, optionID)
		if err != nil {
			return nil, err
		}
		for voterRows.Next() {
			var voterID uuid.UUID
			if err := voterRows.Scan(&voterID); err != nil {
				voterRows.Close()
				return nil, err
			}
			voterIds = append(voterIds, voterID)
		}
		voterRows.Close()

		res.OptionID = optionID
		res.OptionName = optionName
		res.VoteCount = voteCount
		res.VoterIDs = voterIds

		results = append(results, res)
	}

	return results, nil
}

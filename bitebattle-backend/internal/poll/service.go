package poll

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

// Create a new poll for a group by a user
func (s *Service) CreatePoll(name string, createdBy uuid.UUID) (*Poll, error) {
	id := uuid.New()
	now := time.Now()
	inviteCode := utils.GenerateRandomString(8) // Make sure you have this util

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

	// Insert the creator into the poll participants
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

func (s *Service) GetPolls(userID uuid.UUID) ([]PollSummary, error) {
	rows, err := s.DB.Query(`
		SELECT p.id, p.name, p.invite_code,
			CASE 
				WHEN p.created_by = $1 THEN 'owner'
				ELSE 'member'
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

	var polls []PollSummary
	for rows.Next() {
		var p PollSummary
		if err := rows.Scan(&p.ID, &p.Name, &p.InviteCode, &p.Role); err != nil {
			return nil, err
		}
		polls = append(polls, p)
	}

	return polls, nil
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

// Get a poll by its ID
func (s *Service) GetPoll(pollID uuid.UUID) (*Poll, error) {
	row := s.DB.QueryRow(`
		SELECT id, name, invite_code, created_by, created_at, updated_at
		FROM polls WHERE id = $1
	`, pollID)
	var poll Poll
	if err := row.Scan(&poll.ID, &poll.Name, &poll.InviteCode, &poll.CreatedBy, &poll.CreatedAt, &poll.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Poll not found
		}
		return nil, err // Other error
	}
	return &poll, nil
}

// Add a restaurant option to the poll
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

// Cast a vote for a given option by a user
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

// Retrieve voting results for a poll
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

	var results []PollResult
	for rows.Next() {
		var res PollResult
		var voterIds []uuid.UUID
		var optionID uuid.UUID
		var optionName string
		var voteCount int

		// Scan option basic info
		if err := rows.Scan(&optionID, &optionName, &voteCount); err != nil {
			return nil, err
		}

		// Query voter ids for this option
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

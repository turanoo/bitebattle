package poll

import (
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

// Create a new poll for a group by a user
func (s *Service) CreatePoll(groupID, createdBy uuid.UUID) (*Poll, error) {
	id := uuid.New()
	now := time.Now()

	_, err := s.DB.Exec(`
		INSERT INTO polls (id, group_id, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, id, groupID, createdBy, now, now)

	if err != nil {
		return nil, err
	}

	return &Poll{
		ID:        id,
		GroupID:   groupID,
		CreatedBy: createdBy,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
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
		if err := rows.Scan(&res.OptionID, &res.OptionName, &res.VoteCount); err != nil {
			return nil, err
		}
		results = append(results, res)
	}

	return results, nil
}

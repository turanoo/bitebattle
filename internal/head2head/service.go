package head2head

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/turanoo/bitebattle/pkg/db"
)

type Service struct {
	DB      *sql.DB
	Matcher *Matcher
}

func NewService(db *sql.DB) *Service {
	return &Service{DB: db, Matcher: NewMatcher(db)}
}

func (s *Service) CreateMatch(inviterID, inviteeID uuid.UUID, categories []string) (*Match, error) {
	id := uuid.New()
	now := time.Now()

	_, err := s.DB.Exec(`
		INSERT INTO head2head_matches (id, inviter_id, invitee_id, status, categories, created_at, updated_at)
		VALUES ($1, $2, $3, 'pending', $4, $5, $6)
	`, id, inviterID, inviteeID, pq.Array(categories), now, now)

	if err != nil {
		return nil, err
	}

	return &Match{
		ID:         id,
		InviterID:  inviterID,
		InviteeID:  inviteeID,
		Status:     "pending",
		Categories: categories,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func (s *Service) AcceptMatch(matchID, userID uuid.UUID) error {
	row := s.DB.QueryRow(`
		SELECT invitee_id FROM head2head_matches WHERE id = $1 AND status = 'pending'
	`, matchID)

	var inviteeID uuid.UUID
	err := db.ScanOne(row, &inviteeID)
	if err != nil {
		return err
	}
	if inviteeID == uuid.Nil || inviteeID != userID {
		return sql.ErrNoRows
	}

	_, err = s.DB.Exec(`
		UPDATE head2head_matches SET status = 'active', updated_at = $2 WHERE id = $1
	`, matchID, time.Now())
	return err
}

func (s *Service) SubmitSwipe(matchID, userID uuid.UUID, restaurantID, restaurantName string, liked bool) (*Swipe, error) {
	id := uuid.New()
	now := time.Now()

	_, err := s.DB.Exec(`
		INSERT INTO head2head_swipes (id, match_id, user_id, restaurant_id, restaurant_name, liked, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, id, matchID, userID, restaurantID, restaurantName, liked, now)

	if err != nil {
		return nil, err
	}

	return &Swipe{
		ID:             id,
		MatchID:        matchID,
		UserID:         userID,
		RestaurantID:   restaurantID,
		RestaurantName: restaurantName,
		Liked:          liked,
		CreatedAt:      now,
	}, nil
}

func (s *Service) GetMutualLikes(matchID uuid.UUID) ([]Swipe, error) {
	return s.Matcher.FindMutualLikes(matchID)
}

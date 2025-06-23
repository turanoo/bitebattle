package head2head

import (
	"database/sql"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/turanoo/bitebattle/pkg/config"
	"github.com/turanoo/bitebattle/pkg/db"
)

type Service struct {
	DB      *sql.DB
	Matcher *Matcher
}

func NewService(db *sql.DB, cfg *config.Config) *Service {
	return &Service{DB: db, Matcher: NewMatcher(db)}
}

func (s *Service) CreateMatch(inviterID, inviteeID string, categories []string) (*Match, error) {
	id := generateRandomID()
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

func (s *Service) AcceptMatch(matchID, userID string) error {
	row := s.DB.QueryRow(`
		SELECT invitee_id FROM head2head_matches WHERE id = $1 AND status = 'pending'
	`, matchID)

	var inviteeID string
	err := db.ScanOne(row, &inviteeID)
	if err != nil {
		return err
	}
	if inviteeID == "" || inviteeID != userID {
		return sql.ErrNoRows
	}

	_, err = s.DB.Exec(`
		UPDATE head2head_matches SET status = 'active', updated_at = $2 WHERE id = $1
	`, matchID, time.Now())
	return err
}

func (s *Service) SubmitSwipe(matchID, userID, restaurantID, restaurantName string, liked bool) (*Swipe, error) {
	id := generateRandomID()
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

func (s *Service) GetMutualLikes(matchID string) ([]Swipe, error) {
	return s.Matcher.FindMutualLikes(matchID)
}

func generateRandomID() string {
	return strings.ReplaceAll(time.Now().Format("20060102150405.000000000"), ".", "")
}

package head2head

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type Matcher struct {
	DB *sql.DB
}

func NewMatcher(db *sql.DB) *Matcher {
	return &Matcher{DB: db}
}

func (m *Matcher) FindMutualLikes(matchID uuid.UUID) ([]Swipe, error) {
	rows, err := m.DB.Query(`
		SELECT restaurant_id, restaurant_name
		FROM head2head_swipes
		WHERE match_id = $1 AND liked = TRUE
		GROUP BY restaurant_id, restaurant_name
		HAVING COUNT(DISTINCT user_id) = 2
	`, matchID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("failed to close rows: %v\n", err)
		}
	}()

	var matches []Swipe
	for rows.Next() {
		var sw Swipe
		if err := rows.Scan(&sw.RestaurantID, &sw.RestaurantName); err != nil {
			return nil, err
		}
		sw.MatchID = matchID
		sw.Liked = true
		matches = append(matches, sw)
	}

	return matches, nil
}

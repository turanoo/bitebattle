package head2head

import (
	"database/sql"

	"github.com/turanoo/bitebattle/pkg/logger"
)

type Matcher struct {
	DB *sql.DB
}

func NewMatcher(db *sql.DB) *Matcher {
	return &Matcher{DB: db}
}

func (m *Matcher) FindMutualLikes(matchID string) ([]Swipe, error) {
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
			logger.Log.WithError(err).Error("failed to close rows")
		}
	}()

	var matches []Swipe
	for rows.Next() {
		var sw Swipe
		if err := rows.Scan(&sw.RestaurantID, &sw.RestaurantName); err != nil {
			return nil, err
		}
		matches = append(matches, sw)
	}
	return matches, nil
}

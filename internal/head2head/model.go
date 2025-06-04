package head2head

import (
	"time"

	"github.com/google/uuid"
)

type Match struct {
	ID         uuid.UUID `json:"id"`
	InviterID  uuid.UUID `json:"inviter_id"`
	InviteeID  uuid.UUID `json:"invitee_id"`
	Status     string    `json:"status"` // pending, active, completed, cancelled
	Categories []string  `json:"categories"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Swipe struct {
	ID             uuid.UUID `json:"id"`
	MatchID        uuid.UUID `json:"match_id"`
	UserID         uuid.UUID `json:"user_id"`
	RestaurantID   string    `json:"restaurant_id"`
	RestaurantName string    `json:"restaurant_name"`
	Liked          bool      `json:"liked"`
	CreatedAt      time.Time `json:"created_at"`
}

package head2head

import (
	"time"

	"github.com/google/uuid"
)

type CreateMatchRequest struct {
	InviteeID  string   `json:"invitee_id" binding:"required,uuid"`
	Categories []string `json:"categories" binding:"required,min=1,dive,required"`
}

type SubmitSwipeRequest struct {
	RestaurantID   string `json:"restaurant_id" binding:"required"`
	RestaurantName string `json:"restaurant_name" binding:"required"`
	Liked          bool   `json:"liked"`
}

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

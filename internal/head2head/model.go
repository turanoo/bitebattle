package head2head

import (
	"time"
)

type CreateMatchRequest struct {
	InviteeID  string   `json:"invitee_id" binding:"required"`
	Categories []string `json:"categories" binding:"required,min=1,dive,required"`
}

type SubmitSwipeRequest struct {
	RestaurantID   string `json:"restaurant_id" binding:"required"`
	RestaurantName string `json:"restaurant_name" binding:"required"`
	Liked          bool   `json:"liked"`
}

type Match struct {
	ID         string    `json:"id"`
	InviterID  string    `json:"inviter_id"`
	InviteeID  string    `json:"invitee_id"`
	Status     string    `json:"status"` // pending, active, completed, cancelled
	Categories []string  `json:"categories"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Swipe struct {
	ID             string    `json:"id"`
	MatchID        string    `json:"match_id"`
	UserID         string    `json:"user_id"`
	RestaurantID   string    `json:"restaurant_id"`
	RestaurantName string    `json:"restaurant_name"`
	Liked          bool      `json:"liked"`
	CreatedAt      time.Time `json:"created_at"`
}

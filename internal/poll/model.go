package poll

import (
	"time"

	"github.com/google/uuid"
)

type CreatePollRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

type JoinPollRequest struct {
	InviteCode string `json:"invite_code" binding:"required,len=8"`
}

type UpdatePollRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

type AddOptionRequest []struct {
	RestaurantID string `json:"restaurant_id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	ImageURL     string `json:"image_url"`
	MenuURL      string `json:"menu_url"`
}

type VoteRequest struct {
	OptionID string `json:"option_id" binding:"required,uuid"`
}

type Poll struct {
	ID         uuid.UUID   `json:"id"`
	Name       string      `json:"name"`
	InviteCode string      `json:"invite_code"`
	Role       string      `json:"role"`
	Members    []uuid.UUID `json:"members"`
	CreatedBy  uuid.UUID   `json:"created_by"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type PollOption struct {
	ID           uuid.UUID `json:"id"`
	PollID       uuid.UUID `json:"poll_id"`
	RestaurantID string    `json:"restaurant_id"`
	Name         string    `json:"name"`
	ImageURL     string    `json:"image_url"`
	MenuURL      string    `json:"menu_url"`
}

type PollVote struct {
	ID        uuid.UUID `json:"id"`
	PollID    uuid.UUID `json:"poll_id"`
	OptionID  uuid.UUID `json:"option_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type PollResult struct {
	OptionID   uuid.UUID   `json:"option_id"`
	OptionName string      `json:"option_name"`
	VoteCount  int         `json:"vote_count"`
	VoterIDs   []uuid.UUID `json:"voter_ids"`
}

type PollSummary struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Role       string    `json:"role"` // "owner" or "member"
	InviteCode string    `json:"invite_code"`
}

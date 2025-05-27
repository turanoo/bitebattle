package poll

import (
	"time"

	"github.com/google/uuid"
)

type Poll struct {
	ID        uuid.UUID
	GroupID   uuid.UUID
	CreatedBy uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PollOption struct {
	ID           uuid.UUID
	PollID       uuid.UUID
	RestaurantID string
	Name         string
	ImageURL     string
	MenuURL      string
}

type PollVote struct {
	ID                uuid.UUID
	PollID            uuid.UUID
	OptionID          uuid.UUID
	UserID            uuid.UUID
	RestaurantPlaceID string    `json:"restaurant_place_id"`
	CreatedAt         time.Time `json:"created_at"`
}

type PollResult struct {
	OptionID   uuid.UUID `json:"option_id"`
	OptionName string    `json:"option_name"`
	VoteCount  int       `json:"vote_count"`
}

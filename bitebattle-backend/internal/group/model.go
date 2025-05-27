package group

import (
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	CreatedBy  uuid.UUID `json:"created_by"`
	InviteCode string    `json:"invite_code"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

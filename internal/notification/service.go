package notification

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	DB *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) Send(userID uuid.UUID, message string) error {
	_, err := s.DB.Exec(`
		INSERT INTO notifications (id, user_id, message, created_at)
		VALUES ($1, $2, $3, $4)
	`, uuid.New(), userID, message, time.Now())
	return err
}

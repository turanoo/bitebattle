package notification

import (
	"database/sql"
	"time"
)

type Service struct {
	DB *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) Send(userID string, message string) error {
	id := generateRandomID() // Use your utils or a helper for random string IDs
	_, err := s.DB.Exec(`
		INSERT INTO notifications (id, user_id, message, created_at)
		VALUES ($1, $2, $3, $4)
	`, id, userID, message, time.Now())
	return err
}

// Helper for random string IDs (replace with your utils if available)
func generateRandomID() string {
	return time.Now().Format("20060102150405.000000000")
}

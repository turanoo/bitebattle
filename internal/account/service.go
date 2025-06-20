package account

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/turanoo/bitebattle/pkg/config"
	"github.com/turanoo/bitebattle/pkg/db"
)

var (
	ErrInvalidPassword = errors.New("current password is incorrect")
	ErrEmailExists     = errors.New("user with this email already exists")
)

type Service struct {
	DB            *sql.DB
	ProfileBucket string
	ObjectUrl     string
}

func NewService(db *sql.DB, cfg *config.Config) *Service {
	return &Service{
		DB:            db,
		ProfileBucket: cfg.GCS.ProfileBucket,
		ObjectUrl:     cfg.GCS.ObjectURL,
	}
}

type UserProfile struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	PhoneNumber   *string   `json:"phone_number,omitempty"`
	ProfilePicURL *string   `json:"profile_pic_url,omitempty"`
	Bio           *string   `json:"bio,omitempty"`
	LastLoginAt   *string   `json:"last_login_at,omitempty"`
}

func (s *Service) GetUserProfile(userID uuid.UUID) (*UserProfile, error) {
	row := s.DB.QueryRow(`SELECT id, name, email, phone_number, profile_pic_url, bio, last_login_at FROM users WHERE id = $1`, userID)

	var profile UserProfile
	err := db.ScanOne(row, &profile.ID, &profile.Name, &profile.Email, &profile.PhoneNumber, &profile.ProfilePicURL, &profile.Bio, &profile.LastLoginAt)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID uuid.UUID, name, email string) error {
	_, err := s.DB.ExecContext(ctx, `
		UPDATE users SET name = $1, email = $2, updated_at = NOW()
		WHERE id = $3
	`, name, email, userID)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrEmailExists
		}
		return err
	}
	return nil
}

func (s *Service) GenerateProfilePicUploadURL(ctx context.Context, userID uuid.UUID) (uploadURL, objectURL string, err error) {
	object := "profile_pics/" + userID.String() + "_" + time.Now().Format("20060102150405") + ".jpg"
	contentType := "image/jpeg"

	url, err := generateSignedUploadURL(ctx, s.ProfileBucket, object, contentType)
	if err != nil {
		return "", "", err
	}
	return url, s.ObjectUrl + s.ProfileBucket + "/" + object, nil
}

func (s *Service) GenerateProfilePicAccessURL(ctx context.Context, userID uuid.UUID) (string, error) {
	profile, err := s.GetUserProfile(userID)
	if err != nil || profile.ProfilePicURL == nil {
		return "", errors.New("profile picture not found")
	}
	bucket, object, err := s.parseGCSURL(*profile.ProfilePicURL)
	if err != nil {
		return "", err
	}
	return generateSignedAccessURL(ctx, bucket, object)
}

func generateSignedUploadURL(ctx context.Context, bucket, object, contentType string) (string, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer func() {
		cerr := client.Close()
		if cerr != nil {
			fmt.Printf("Failed to close storage client: %v\n", cerr)
		}
	}()

	url, err := storage.SignedURL(bucket, object, &storage.SignedURLOptions{
		Method:      "PUT",
		Expires:     time.Now().Add(15 * time.Minute),
		ContentType: contentType,
	})
	return url, err
}

func generateSignedAccessURL(ctx context.Context, bucket, object string) (string, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer func() {
		cerr := client.Close()
		if cerr != nil {
			fmt.Printf("Failed to close storage client: %v\n", cerr)
		}
	}()

	url, err := storage.SignedURL(bucket, object, &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(15 * time.Minute),
	})
	return url, err
}

func (s *Service) parseGCSURL(gcsURL string) (bucket, object string, err error) {
	if !strings.HasPrefix(gcsURL, s.ObjectUrl) {
		return "", "", fmt.Errorf("invalid GCS URL")
	}
	trimmed := strings.TrimPrefix(gcsURL, s.ObjectUrl)
	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid GCS URL")
	}
	return parts[0], parts[1], nil
}

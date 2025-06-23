package account

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/turanoo/bitebattle/pkg/config"
	"github.com/turanoo/bitebattle/pkg/logger"
)

var (
	ErrInvalidPassword = errors.New("current password is incorrect")
	ErrEmailExists     = errors.New("user with this email already exists")
)

type Service struct {
	ProfileBucket string
	ObjectUrl     string
	Auth0         config.Auth0
}

func NewService(_ interface{}, cfg *config.Config) *Service {
	return &Service{
		ProfileBucket: cfg.GCS.ProfileBucket,
		ObjectUrl:     cfg.GCS.ObjectURL,
		Auth0:         cfg.Auth0,
	}
}

type UserProfile struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	PhoneNumber   *string `json:"phone_number,omitempty"`
	ProfilePicURL *string `json:"profile_pic_url,omitempty"`
	Bio           *string `json:"bio,omitempty"`
	LastLoginAt   *string `json:"last_login_at,omitempty"`
}

type auth0User struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

func (s *Service) getManagementToken(ctx context.Context) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", s.Auth0.ClientID)
	data.Set("client_secret", s.Auth0.ClientSecret)
	data.Set("audience", s.Auth0.ManagementAPI)

	req, err := http.NewRequestWithContext(ctx, "POST", s.Auth0.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get management token: %s", string(body))
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.AccessToken, nil
}

func (s *Service) GetUserProfile(userID string) (*UserProfile, error) {
	ctx := context.Background()
	token, err := s.getManagementToken(ctx)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/users/%s", s.Auth0.ManagementAPI, userID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch user profile: %s", string(body))
	}

	var user auth0User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &UserProfile{
		ID:    user.UserID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID, name, email string) error {
	token, err := s.getManagementToken(ctx)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/users/%s", s.Auth0.ManagementAPI, userID)
	payload := map[string]interface{}{
		"name":  name,
		"email": email,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to update user profile: %s", string(respBody))
	}
	return nil
}

func (s *Service) GenerateProfilePicUploadURL(ctx context.Context, userID string) (uploadURL, objectURL string, err error) {
	object := "profile_pics/" + userID + "_" + time.Now().Format("20060102150405") + ".jpg"
	contentType := "image/jpeg"

	url, err := generateSignedUploadURL(ctx, s.ProfileBucket, object, contentType)
	if err != nil {
		return "", "", err
	}
	return url, s.ObjectUrl + s.ProfileBucket + "/" + object, nil
}

func (s *Service) GenerateProfilePicAccessURL(ctx context.Context, userID string) (string, error) {
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
			logger.Log.WithError(cerr).Error("Failed to close storage client")
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
			logger.Log.WithError(cerr).Error("Failed to close storage client")
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

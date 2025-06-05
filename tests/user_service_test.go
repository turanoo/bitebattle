package tests

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/turanoo/bitebattle/internal/user"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	_, err = db.Exec(`
		CREATE TABLE users (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL,
			name TEXT NOT NULL,
			password_hash TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);
	`)
	if err != nil {
		t.Fatalf("failed to create users table: %v", err)
	}
	return db
}

func TestCreateAndGetUser(t *testing.T) {
	db := setupTestDB(t)
	service := user.NewService(db)
	ctx := context.Background()

	u := &user.User{
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: "hashedpassword",
	}
	created, err := service.CreateUser(ctx, u)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if created.ID == "" {
		t.Error("expected user ID to be set")
	}
	fetched, err := service.GetUserByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}
	if fetched.Email != u.Email {
		t.Errorf("expected email %s, got %s", u.Email, fetched.Email)
	}
}

func TestGetUserByEmail(t *testing.T) {
	db := setupTestDB(t)
	service := user.NewService(db)
	ctx := context.Background()

	u := &user.User{
		Email:        "emailtest@example.com",
		Name:         "Email Test",
		PasswordHash: "hashedpassword",
	}
	_, err := service.CreateUser(ctx, u)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	fetched, err := service.GetUserByEmail(ctx, u.Email)
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}
	if fetched.Name != u.Name {
		t.Errorf("expected name %s, got %s", u.Name, fetched.Name)
	}
}

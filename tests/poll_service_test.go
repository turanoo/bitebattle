package tests

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/turanoo/bitebattle/internal/poll"
)

func setupPollTestDB(t *testing.T) *sql.DB {
	dbName := fmt.Sprintf("file:polltest_%d?mode=memory&cache=shared&_foreign_keys=on", time.Now().UnixNano())
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	_, err = db.Exec(`CREATE TABLE users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		password_hash TEXT NOT NULL,
		created_at DATETIME,
		updated_at DATETIME
	);`)
	if err != nil {
		t.Fatalf("failed to create users table: %v", err)
	}
	_, err = db.Exec(`CREATE TABLE polls (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		invite_code TEXT NOT NULL,
		created_by TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);`)
	if err != nil {
		t.Fatalf("failed to create polls table: %v", err)
	}
	_, err = db.Exec(`CREATE TABLE polls_members (
		id TEXT PRIMARY KEY,
		poll_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		joined_at DATETIME,
		UNIQUE (poll_id, user_id)
	);`)
	if err != nil {
		t.Fatalf("failed to create polls_members table: %v", err)
	}
	_, err = db.Exec(`CREATE TABLE poll_options (
		id TEXT PRIMARY KEY,
		poll_id TEXT NOT NULL,
		restaurant_id TEXT NOT NULL,
		name TEXT NOT NULL,
		image_url TEXT,
		menu_url TEXT
	);`)
	if err != nil {
		t.Fatalf("failed to create poll_options table: %v", err)
	}
	_, err = db.Exec(`CREATE TABLE poll_votes (
		id TEXT PRIMARY KEY,
		poll_id TEXT NOT NULL,
		option_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		created_at DATETIME
	);`)
	if err != nil {
		t.Fatalf("failed to create poll_votes table: %v", err)
	}
	return db
}

func TestCreatePoll(t *testing.T) {
	db := setupPollTestDB(t)
	service := poll.NewService(db)
	userID := uuid.New()
	p, err := service.CreatePoll("Test Poll", userID)
	if err != nil {
		t.Fatalf("CreatePoll failed: %v", err)
	}
	if p.Name != "Test Poll" {
		t.Errorf("expected poll name 'Test Poll', got %s", p.Name)
	}
}

func TestGetPolls(t *testing.T) {
	db := setupPollTestDB(t)
	service := poll.NewService(db)
	userID := uuid.New()
	_, err := service.CreatePoll("Poll1", userID)
	if err != nil {
		t.Fatalf("CreatePoll failed: %v", err)
	}
	polls, err := service.GetPolls(userID)
	if err != nil {
		t.Fatalf("GetPolls failed: %v", err)
	}
	if len(polls) == 0 {
		t.Error("expected at least one poll")
	}
}

func TestDeletePoll(t *testing.T) {
	db := setupPollTestDB(t)
	service := poll.NewService(db)
	userID := uuid.New()
	p, err := service.CreatePoll("PollToDelete", userID)
	if err != nil {
		t.Fatalf("CreatePoll failed: %v", err)
	}
	err = service.DeletePoll(p.ID)
	if err != nil {
		t.Fatalf("DeletePoll failed: %v", err)
	}
	polls, err := service.GetPolls(userID)
	if err != nil {
		t.Fatalf("GetPolls failed: %v", err)
	}
	for _, poll := range polls {
		if poll.ID == p.ID {
			t.Error("poll was not deleted")
		}
	}
}

package tests

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/turanoo/bitebattle/internal/account"
	"github.com/turanoo/bitebattle/pkg/utils"
)

func setupAccountTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	_, err = db.Exec(`
		CREATE TABLE users (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL,
			name TEXT NOT NULL,
			password_hash TEXT NOT NULL
		);
	`)
	if err != nil {
		t.Fatalf("failed to create users table: %v", err)
	}
	return db
}

func TestGetUserProfile(t *testing.T) {
	db := setupAccountTestDB(t)
	service := account.NewService(db)
	id := uuid.New()
	_, err := db.Exec(`INSERT INTO users (id, email, name, password_hash) VALUES (?, ?, ?, ?)`, id.String(), "test@a.com", "Test", "hash")
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}
	profile, err := service.GetUserProfile(id)
	if err != nil {
		t.Fatalf("GetUserProfile failed: %v", err)
	}
	if profile.Name != "Test" {
		t.Errorf("expected name 'Test', got %s", profile.Name)
	}
}

func TestUpdateUserProfile_NameAndEmail(t *testing.T) {
	db := setupAccountTestDB(t)
	service := account.NewService(db)
	id := uuid.New()
	_, err := db.Exec(`INSERT INTO users (id, email, name, password_hash) VALUES (?, ?, ?, ?)`, id.String(), "old@a.com", "OldName", "hash")
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}
	newName := "NewName"
	newEmail := "new@a.com"
	err = service.UpdateUserProfile(id, &newName, &newEmail, nil, nil)
	if err != nil {
		t.Fatalf("UpdateUserProfile failed: %v", err)
	}
	row := db.QueryRow(`SELECT name, email FROM users WHERE id = ?`, id.String())
	var gotName, gotEmail string
	err = row.Scan(&gotName, &gotEmail)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if gotName != newName || gotEmail != newEmail {
		t.Errorf("expected %s/%s, got %s/%s", newName, newEmail, gotName, gotEmail)
	}
}

func TestUpdateUserProfile_Password(t *testing.T) {
	db := setupAccountTestDB(t)
	service := account.NewService(db)
	id := uuid.New()
	oldHash, _ := utils.HashPassword("oldpass")
	_, err := db.Exec(`INSERT INTO users (id, email, name, password_hash) VALUES (?, ?, ?, ?)`, id.String(), "a@b.com", "Name", oldHash)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}
	cur := "oldpass"
	newp := "newpass"
	err = service.UpdateUserProfile(id, nil, nil, &cur, &newp)
	if err != nil {
		t.Fatalf("UpdateUserProfile failed: %v", err)
	}
	row := db.QueryRow(`SELECT password_hash FROM users WHERE id = ?`, id.String())
	var gotHash string
	err = row.Scan(&gotHash)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if err := utils.CheckPasswordHash(gotHash, newp); err != nil {
		t.Errorf("password hash not updated correctly")
	}
}

func TestUpdateUserProfile_Errors(t *testing.T) {
	db := setupAccountTestDB(t)
	service := account.NewService(db)
	id := uuid.New()
	_, err := db.Exec(`INSERT INTO users (id, email, name, password_hash) VALUES (?, ?, ?, ?)`, id.String(), "a@b.com", "Name", "hash")
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}
	// No fields
	err = service.UpdateUserProfile(id, nil, nil, nil, nil)
	if err == nil {
		t.Error("expected error for no fields")
	}
	// Only one password field
	cur := "x"
	err = service.UpdateUserProfile(id, nil, nil, &cur, nil)
	if err == nil {
		t.Error("expected error for only one password field")
	}
}

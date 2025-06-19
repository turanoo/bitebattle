package user

import "time"

type User struct {
	ID            string     `db:"id" json:"id"`
	Email         string     `db:"email" json:"email"`
	Name          string     `db:"name" json:"name"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
	PasswordHash  string     `db:"password_hash" json:"-"`
	PhoneNumber   *string    `db:"phone_number" json:"phone_number,omitempty"`
	ProfilePicURL *string    `db:"profile_pic_url" json:"profile_pic_url,omitempty"`
	Bio           *string    `db:"bio" json:"bio,omitempty"`
	LastLoginAt   *time.Time `db:"last_login_at" json:"last_login_at,omitempty"`
}

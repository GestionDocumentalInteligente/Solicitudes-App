package datamod

import (
	"time"
)

type User struct {
	ID                   int64      `json:"id" db:"id"`
	PersonID             *int64     `json:"person_id,omitempty" db:"person_id"`
	EmailValidated       bool       `json:"email_validated" db:"email_validated"`
	AcceptsNotifications bool       `json:"accepts_notifications" db:"accepts_notifications"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	LoggedAt             *time.Time `json:"logged_at,omitempty" db:"logged_at"`
	UpdatedAt            *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt            *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

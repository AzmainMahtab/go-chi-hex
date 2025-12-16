// Package domain
// Domain holds the source of truth for our data
package domain

import "time"

type User struct {
	ID         int
	UserName   string
	Phone      string
	Password   string
	UserStatus string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

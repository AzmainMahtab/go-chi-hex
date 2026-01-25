// Package domain
// Domain holds the source of truth for our data
package domain

import "time"

type User struct {
	ID         int        `db:"id"`
	UUID       string     `db:"uuid"`
	UserName   string     `db:"user_name"`
	Email      string     `db:"email"`
	Phone      string     `db:"phone"`
	Password   string     `db:"password"`
	UserStatus string     `db:"user_status"`
	UserRole   string     `db:"user_role"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
}

type UserFilter struct {
	UserName    string
	Email       string
	Phone       string
	ShowDeleted bool
	UserStatus  string
	Limit       int
	Offset      int
}

type UserUpdate struct {
	ID       string
	UserName *string
	Email    *string
	Phone    *string
	Status   *string
}

// Package domain
// this one is the auth domain
package domain

type AuthLogin struct {
	Email    string
	Password string
}

type UserClaims struct {
	UserID   string
	Email    string
	Role     string
	IssuedAt int64
	Expires  int64
}

type Tokenpair struct {
	AccessToken string
	RefreshToke string
}

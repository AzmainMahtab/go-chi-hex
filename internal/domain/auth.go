// Package domain
// this one is the auth domain
package domain

type AuithLogin struct {
	Email    string
	Password string
}

type UserClaims struct {
	UserID int
	Email  string
	Role   string
}

type Tokenpair struct {
	AccessToken string
	RefreshToke string
}

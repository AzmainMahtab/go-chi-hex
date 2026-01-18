// Package secure
// this contains JWT related stuff
package secure

import (
	"crypto/ecdsa"
	"time"
	//
	// "github.com/AzmainMahtab/go-chi-hex/internal/domain"
	// "github.com/golang-jwt/jwt/v4"
)

type jwtAdapter struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	accessTTL  time.Duration
	refreshTTL time.Duration
	issuer     string
}

func NewJWT(
	priv *ecdsa.PrivateKey,
	pub *ecdsa.PublicKey,
	aTTL time.Duration,
	rTTL time.Duration,
	iss string,
) *jwtAdapter {
	return &jwtAdapter{
		privateKey: priv,
		publicKey:  pub,
		accessTTL:  aTTL,
		refreshTTL: rTTL,
		issuer:     iss,
	}
}

// func (j *jwtAdapter) GenerateTokenPair (user *domain.User) (domain.Tokenpair, error) {
// 	return nil, nil
// }

// func (j *jwtAdapter) signToken (user *domain.User, ttl time.Duration) (string, error) {
// 	claim := jwt.MapClaims{
// 		"sub": user.ID,
// 		"email": user.Email,
// 		"role": user.Email
// 	}
// }

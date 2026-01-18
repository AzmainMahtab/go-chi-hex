// Package secure
// this contains JWT related stuff
package secure

import (
	"crypto/ecdsa"
	"time"

	"github.com/AzmainMahtab/go-chi-hex/internal/domain"
	"github.com/golang-jwt/jwt/v5"
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

func (j *jwtAdapter) GenerateTokenPair(user *domain.User) (domain.Tokenpair, error) {
	accToken, err := j.signToken(user, j.accessTTL)
	if err != nil {
		return domain.Tokenpair{}, err
	}

	refToken, err := j.signToken(user, j.refreshTTL)
	if err != nil {
		return domain.Tokenpair{}, err
	}

	return domain.Tokenpair{
		AccessToken: accToken,
		RefreshToke: refToken,
	}, err
}

func (j *jwtAdapter) signToken(u *domain.User, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":   u.ID,
		"email": u.Email,
		"role":  u.UserRole,
		"iss":   j.issuer,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	signed, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", err
	}

	return signed, nil
}

// Package secure
// this contains JWT related stuff
package secure

import (
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/AzmainMahtab/go-chi-hex/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type JWTAdapter struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	AccessTTL  time.Duration
	RefreshTTL time.Duration
	Issuer     string
}

func NewJWT(
	priv *ecdsa.PrivateKey,
	pub *ecdsa.PublicKey,
	aTTL time.Duration,
	rTTL time.Duration,
	iss string,
) *JWTAdapter {
	return &JWTAdapter{
		PrivateKey: priv,
		PublicKey:  pub,
		AccessTTL:  aTTL,
		RefreshTTL: rTTL,
		Issuer:     iss,
	}
}

func (j *JWTAdapter) GenerateTokenPair(user *domain.User) (domain.Tokenpair, error) {
	accToken, err := j.signToken(user, j.AccessTTL)
	if err != nil {
		return domain.Tokenpair{}, err
	}

	refToken, err := j.signToken(user, j.RefreshTTL)
	if err != nil {
		return domain.Tokenpair{}, err
	}

	return domain.Tokenpair{
		AccessToken: accToken,
		RefreshToke: refToken,
	}, err
}

func (j *JWTAdapter) VerifyToken(tokenStr string) (domain.UserClaims, error) {
	//  Parse the token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token is using the ECDSA algorithm we expect
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Use the PUBLIC key to verify the signature
		return j.PublicKey, nil
	})

	//  Handle parsing errors or invalid tokens
	if err != nil || !token.Valid {
		return domain.UserClaims{}, &domain.AppError{
			Code:    domain.CodeInvalidToken,
			Message: "Bad token",
			Err:     err,
		}
	}

	//  Extract Claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return domain.UserClaims{}, &domain.AppError{
			Code:    domain.CodeInvalidToken,
			Message: "Bad token",
			Err:     err,
		}
	}

	//  Map the map[string]any back to your clean Domain struct
	// IMPORTANT jwt.MapClaims stores numbers as float64
	return domain.UserClaims{
		UserID:   claims["sub"].(string),
		Email:    claims["email"].(string),
		Role:     claims["role"].(string),
		IssuedAt: int64(claims["iat"].(float64)),
		Expires:  int64(claims["exp"].(float64)),
	}, nil
}

func (j *JWTAdapter) signToken(u *domain.User, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":   u.UUID,
		"email": u.Email,
		"role":  u.UserRole,
		"iss":   j.Issuer,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	signed, err := token.SignedString(j.PrivateKey)
	if err != nil {
		return "", err
	}

	return signed, nil
}

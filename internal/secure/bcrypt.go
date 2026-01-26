// Package secure
// this one contains the bcrypt hasing logic
package secure

import (
	"log/slog"

	"github.com/AzmainMahtab/go-chi-hex/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	Cost int
}

func NewBcryptHasher(cost int) *BcryptHasher {
	return &BcryptHasher{
		Cost: cost,
	}
}

func (b *BcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), b.Cost)
	if err != nil {
		slog.Error("Error with password hasing:")
		return "", &domain.AppError{
			Code:    domain.CodeInternal,
			Message: "Something went wrong",
			Err:     err,
		}
	}

	return string(bytes), nil
}

func (b *BcryptHasher) Compare(hashed string, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err == nil
}

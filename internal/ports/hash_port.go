// Package ports
// This port contains pasword hashing and other hashing
package ports

// Adapter is in internal/secure directory
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash string, plain string) bool
}

package hash

import (
	"crypto/sha1"
	"fmt"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

// SHA1PasswordHasher uses SHA1 to hash passwords with provided salt
type SHA1PasswordHasher struct {
	salt string
}

func NewSHA1PasswordHasher(salt string) *SHA1PasswordHasher {
	return &SHA1PasswordHasher{salt: salt}
}

// Hash creates SHA1 hash of given password
func (h *SHA1PasswordHasher) Hash(password string) (string, error) {
	hash := sha1.New()

	_, err := hash.Write([]byte(password))

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt))), nil
}

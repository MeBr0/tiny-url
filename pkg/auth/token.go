package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// TokenManager provides token issuing and decoding
type TokenManager interface {
	Issue(subject string, ttl time.Duration) (string, error)
	Decode(token string) (string, error)
}

type JWTManager struct {
	signingKey string
}

func NewJWTManager(signingKey string) (*JWTManager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &JWTManager{signingKey: signingKey}, nil
}

func (m *JWTManager) Issue(subject string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(ttl).Unix(),
		Subject:   subject,
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *JWTManager) Decode(token string) (string, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error get user claims from t")
	}

	return claims["sub"].(string), nil
}

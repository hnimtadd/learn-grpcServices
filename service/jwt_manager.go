package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

var ErrUnexpectedSigningMethod = errors.New("Unexpected signing method")

type UserClaim struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{secretKey: secretKey, tokenDuration: tokenDuration}
}

func (manager *JWTManager) Generate(user *User) (string, error) {
	claims := &UserClaim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(manager.tokenDuration).Unix(),
		},
		Username: user.Username,
		Role:     user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(manager.secretKey))
	return ss, err
}

func (manager *JWTManager) Verify(accessToken string) (*UserClaim, error) {
	token, err := jwt.ParseWithClaims(accessToken, &UserClaim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}
		return []byte(manager.secretKey), nil

	})
	if err != nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*UserClaim)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}

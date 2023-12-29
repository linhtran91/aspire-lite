package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	CustomerID int64 `json:"customer_id"`
	jwt.RegisteredClaims
}

type jwtTokenBuilder struct {
	secretKey []byte
	duration  time.Duration
}

func NewJWTTokenBuilder(secret string, duration time.Duration) *jwtTokenBuilder {
	return &jwtTokenBuilder{
		secretKey: []byte(secret),
		duration:  duration,
	}
}

func (b *jwtTokenBuilder) Encode(customerID int64) (string, error) {
	now := time.Now()
	claims := TokenClaims{
		CustomerID: customerID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(b.duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(b.secretKey)
}

func (b *jwtTokenBuilder) Decode(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return b.secretKey, nil
	})
	if err != nil {
		return -1, err
	}
	if claims, ok := token.Claims.(*TokenClaims); ok {
		return claims.CustomerID, nil
	}
	return -1, errors.New("unknown claims type, cannot proceed")
}

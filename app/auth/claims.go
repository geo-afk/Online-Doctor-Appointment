package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	UserId   int32  `json:"user_name"`
	UserRole string `json:"user_role"`

	jwt.RegisteredClaims
}

func NewUserClaims(userId int32, userRole string, duration time.Duration) (*UserClaims, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("error generating token id: %w", err)
	}

	return &UserClaims{
		UserId:   userId,
		UserRole: userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenId.String(),
			Subject:   userRole,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}

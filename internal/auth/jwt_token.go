package auth

import (
	"strconv"
	"time"

	"github.com/crazydw4rf/oil-bank-backend/internal/entity"
	"github.com/crazydw4rf/oil-bank-backend/internal/services/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	CollectorID string `json:"collector_id"`
	jwt.RegisteredClaims
}

func GenerateToken(user *entity.UserWithCollector, cfg *config.Config) (*entity.UserWithCollector, error) {
	at, err := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		CollectorID: strconv.FormatInt(user.CollectorId, 10),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(user.Id, 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}).
		SignedString([]byte(cfg.JWT_ACCESS_TOKEN_SECRET))
	if err != nil {
		return nil, err
	}

	rt, err := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		CollectorID: strconv.FormatInt(user.CollectorId, 10),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(user.Id, 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		},
	}).
		SignedString([]byte(cfg.JWT_REFRESH_TOKEN_SECRET))
	if err != nil {
		return nil, err
	}

	user.AccessToken = at
	user.RefreshToken = rt

	return user, nil
}

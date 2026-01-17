package middleware

import (
	"github.com/crazydw4rf/oil-bank-backend/internal/auth"
	"github.com/crazydw4rf/oil-bank-backend/internal/constants"
	"github.com/crazydw4rf/oil-bank-backend/internal/delivery/http/response"
	. "github.com/crazydw4rf/oil-bank-backend/internal/delivery/http/response"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func (m HTTPMiddleware) Verify(c *fiber.Ctx) error {
	var tokenString string

	// ambil token dari cookie
	tokenString = c.Cookies("token")
	if tokenString == "" {
		return NewHTTPErrorSimple(c, fiber.StatusUnauthorized, "Missing token")
	}
	// else if tokenHeader := c.Get(config.ACCESS_TOKEN_HEADER_NAME); tokenHeader != "" {
	// 	tokenString = strings.TrimPrefix(tokenHeader, "Bearer ")
	// 	if tokenString == tokenHeader {
	// 		return NewHTTPErrorSimple(c, fiber.StatusBadRequest, "Invalid token")
	// 	}
	// } else {
	// 	return NewHTTPErrorSimple(c, fiber.StatusUnauthorized, "Missing token")
	// }

	claims := new(auth.JWTClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if t.Header["alg"] != "HS256" {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(m.cfg.JWT_ACCESS_TOKEN_SECRET), nil
	})
	// FIXME: buat centralized error handling
	if err != nil {
		switch err {
		case jwt.ErrTokenExpired:
			return NewHTTPErrorSimple(c, fiber.StatusUnauthorized, "Token expired")
		case jwt.ErrSignatureInvalid:
		case jwt.ErrTokenMalformed:
			return NewHTTPErrorSimple(c, fiber.StatusBadRequest, "Token has invalid signature or malformed")
		default:
			return NewHTTPErrorSimple(c, fiber.StatusInternalServerError, "Unknown error")
		}
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		return response.NewHTTPErrorSimple(c, fiber.StatusInternalServerError, "Failed to get subject from token")
	}

	c.Locals(constants.UserIdKey, sub)
	c.Locals(constants.CollectorIdKey, claims.CollectorID)

	return c.Next()
}

func (m HTTPMiddleware) VerifyRefreshToken(c *fiber.Ctx) error {
	// TODO: implementasi logika refresh token
	return c.Next()
}

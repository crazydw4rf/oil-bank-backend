package middleware

import "github.com/crazydw4rf/oil-bank-backend/internal/services/config"

type HTTPMiddleware struct {
	cfg *config.Config
}

func NewHTTPMiddleware(cfg *config.Config) HTTPMiddleware {
	return HTTPMiddleware{cfg}
}

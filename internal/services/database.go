package services

import (
	"fmt"

	"github.com/crazydw4rf/oil-bank-backend/internal/services/config"
	"github.com/gofiber/fiber/v2/log"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type DatabaseService struct {
	*sqlx.DB
}

func NewDatabaseService(cfg *config.Config) (DatabaseService, error) {
	db, err := sqlx.Connect("pgx", cfg.DATABASE_URL)
	if err != nil {
		log.Errorf("%#v\n", err)
		return DatabaseService{}, fmt.Errorf("failed to connect to database: %w", err)
	}

	return DatabaseService{db}, nil
}

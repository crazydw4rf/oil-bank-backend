package services

import (
	"github.com/crazydw4rf/oil-bank-backend/internal/services/config"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func NewFiberService() *fiber.App {
	isDevMode := config.APP_ENV == "development"
	app := fiber.New(fiber.Config{
		EnablePrintRoutes:     isDevMode,
		AppName:               "UCOB! API",
		DisableStartupMessage: !isDevMode,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
	})

	return app
}

package main

import (
	"context"
	"fmt"

	"github.com/crazydw4rf/oil-bank-backend/internal/delivery/http/controller"
	"github.com/crazydw4rf/oil-bank-backend/internal/delivery/http/middleware"
	"github.com/crazydw4rf/oil-bank-backend/internal/repository"
	"github.com/crazydw4rf/oil-bank-backend/internal/services"
	"github.com/crazydw4rf/oil-bank-backend/internal/services/config"
	"github.com/crazydw4rf/oil-bank-backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

func start(lc fx.Lifecycle, app *fiber.App, cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Starting app...")
			addr := fmt.Sprintf("%s:%d", cfg.APP_HOST, cfg.APP_PORT)
			go app.Listen(addr)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Stopping app...")
			err := app.ShutdownWithContext(ctx)
			return err
		},
	})
}

func publicRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"app_name": "oil-bank-backend", "version": config.APP_VERSION})
	})
}

func main() {
	app := fx.New(
		fx.Provide(config.InitConfig, services.NewFiberService, services.NewDatabaseService),
		fx.Provide(middleware.NewHTTPMiddleware),
		fx.Provide(repository.NewUserRepository, usecase.NewUserUsecase, controller.NewUserController),
		fx.Provide(repository.NewTransactionRepository, usecase.NewTransactionUsecase, controller.NewTransactionController),
		fx.Provide(repository.NewReportRepository, usecase.NewReportUsecase, controller.NewReportController),
		fx.Provide(repository.NewOilRepository, usecase.NewOilUsecase, controller.NewOilController),
		fx.Invoke(publicRoutes, controller.SetupUserRouter, controller.SetupOilRouter, controller.SetupTransactionRouter, controller.SetupReportRouter),
		fx.Invoke(start),
	)

	app.Run()
}

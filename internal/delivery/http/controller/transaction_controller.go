package controller

import (
	"log"

	"github.com/crazydw4rf/oil-bank-backend/internal/delivery/http/middleware"
	. "github.com/crazydw4rf/oil-bank-backend/internal/delivery/http/response"
	"github.com/crazydw4rf/oil-bank-backend/internal/dto"
	"github.com/crazydw4rf/oil-bank-backend/internal/services/config"
	"github.com/crazydw4rf/oil-bank-backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

const (
	BASE_TRANSACTION_PATH = config.BASE_API_HTTP_PATH + "/transactions"
	TRANSACTION_CREATE    = BASE_TRANSACTION_PATH
	TRANSACTION_UPDATE    = BASE_TRANSACTION_PATH + "/:id"
)

type TransactionController struct {
	transactionUsecase usecase.ITransactionUsecase
}

func NewTransactionController(transactionUsecase usecase.ITransactionUsecase) TransactionController {
	return TransactionController{transactionUsecase}
}

func (tc TransactionController) CreateTransaction(c *fiber.Ctx) error {
	req := new(dto.TransactionCreateDto)
	if err := c.BodyParser(req); err != nil {
		return NewHTTPErrorSimple(c, fiber.StatusBadRequest, "Invalid request body", true)
	}

	var res = CollectorIdExtractor(c)
	if res.IsError() {
		return NewHTTPErrorSimple(c, fiber.StatusBadRequest, "Collector ID not found", true)
	}
	collectorId := res.Value()

	ctx := c.Context()

	result := tc.transactionUsecase.CreateTransaction(ctx, collectorId, req)
	if result.IsError() {
		if err := result.ExpectedError(); err != nil {
			return NewHTTPError(c, err)
		}

		log.Println(result)
		return NewHTTPErrorSimple(c, fiber.StatusInternalServerError, "Failed to create transaction", true)
	}

	return NewHTTPResponse(c, fiber.StatusCreated, result.Value())
}

func (tc TransactionController) UpdateTransaction(c *fiber.Ctx) error {
	// Parse transaction ID from URL parameter
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return NewHTTPErrorSimple(c, fiber.StatusBadRequest, "Invalid transaction ID", true)
	}

	// Parse request body
	req := new(dto.UpdateTransactionDto)
	if err := c.BodyParser(req); err != nil {
		return NewHTTPErrorSimple(c, fiber.StatusBadRequest, "Invalid request body", true)
	}

	// Validate transaction type
	if req.TransactionType != dto.TRANSACTION_SELL && req.TransactionType != dto.TRANSACTION_BUY {
		return NewHTTPErrorSimple(c, fiber.StatusBadRequest, "Invalid transaction type. Must be SELL or BUY", true)
	}

	// Extract collector ID from context
	var res = CollectorIdExtractor(c)
	if res.IsError() {
		return NewHTTPErrorSimple(c, fiber.StatusBadRequest, "Collector ID not found", true)
	}
	collectorId := res.Value()

	ctx := c.Context()

	result := tc.transactionUsecase.UpdateTransaction(ctx, collectorId, int64(id), req)
	if result.IsError() {
		if err := result.ExpectedError(); err != nil {
			return NewHTTPError(c, err)
		}

		log.Println(result)
		return NewHTTPErrorSimple(c, fiber.StatusInternalServerError, "Failed to update transaction", true)
	}

	return NewHTTPResponse(c, fiber.StatusOK, result.Value())
}

func SetupTransactionRouter(app *fiber.App, ctrl TransactionController, mw middleware.HTTPMiddleware) {
	app.Post(TRANSACTION_CREATE, mw.Verify, ctrl.CreateTransaction)
	app.Patch(TRANSACTION_UPDATE, mw.Verify, ctrl.UpdateTransaction)
}

package controller

import (
	"log"
	"strconv"

	"github.com/crazydw4rf/oil-bank-backend/internal/delivery/http/middleware"
	. "github.com/crazydw4rf/oil-bank-backend/internal/delivery/http/response"
	"github.com/crazydw4rf/oil-bank-backend/internal/services/config"
	"github.com/crazydw4rf/oil-bank-backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

const (
	BASE_OIL_PATH = config.BASE_API_HTTP_PATH + "/oil"
	OIL_GETMANY   = "/"
	OIL_GET       = "/:id"
	OIL_UPDATE    = "/:id"
	OIL_DELETE    = "/:id"
)

type OilController struct {
	oilUsecase usecase.IOilUsecase
}

func NewOilController(oilUsecase usecase.IOilUsecase) OilController {
	return OilController{oilUsecase}
}

func (oc OilController) GetOil(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return NewHTTPErrorSimple(c, fiber.StatusBadRequest, "Invalid oil ID", true)
	}

	ctx := c.Context()

	result := oc.oilUsecase.GetOil(ctx, id)
	if result.IsError() {
		if err := result.ExpectedError(); err != nil {
			return NewHTTPError(c, err)
		}

		log.Println(result)
		return NewHTTPErrorSimple(c, fiber.StatusInternalServerError, "Failed to get oil record", true)
	}

	return NewHTTPResponse(c, fiber.StatusOK, result.Value())
}

func (oc OilController) GetOilByCollectorId(c *fiber.Ctx) error {
	res := CollectorIdExtractor(c)
	if res.IsError() {
		return NewHTTPErrorSimple(c, fiber.StatusBadRequest, "Invalid collector ID", true)
	}
	collectorId := res.Value()

	ctx := c.Context()

	result := oc.oilUsecase.GetOilByCollectorId(ctx, collectorId)
	if result.IsError() {
		if err := result.ExpectedError(); err != nil {
			return NewHTTPError(c, err)
		}

		log.Println(result)
		return NewHTTPErrorSimple(c, fiber.StatusInternalServerError, "Failed to get oil inventory", true)
	}

	return NewHTTPResponse(c, fiber.StatusOK, result.Value())
}

func (oc OilController) UpdateOil(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return NewHTTPErrorSimple(c, fiber.StatusBadRequest, "Invalid oil ID", true)
	}

	req := struct {
		TotalVolume float64 `json:"total_volume" validate:"gte=0"`
	}{}

	if err := c.BodyParser(&req); err != nil {
		return NewHTTPErrorSimple(c, fiber.StatusBadRequest, "Invalid request body", true)
	}

	ctx := c.Context()

	result := oc.oilUsecase.UpdateOil(ctx, id, req.TotalVolume)
	if result.IsError() {
		if err := result.ExpectedError(); err != nil {
			return NewHTTPError(c, err)
		}

		log.Println(result)
		return NewHTTPErrorSimple(c, fiber.StatusInternalServerError, "Failed to update oil record", true)
	}

	return NewHTTPResponse(c, fiber.StatusOK, result.Value())
}

func (oc OilController) DeleteOil(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return NewHTTPErrorSimple(c, fiber.StatusBadRequest, "Invalid oil ID", true)
	}

	ctx := c.Context()

	result := oc.oilUsecase.DeleteOil(ctx, id)
	if result.IsError() {
		if err := result.ExpectedError(); err != nil {
			return NewHTTPError(c, err)
		}

		log.Println(result)
		return NewHTTPErrorSimple(c, fiber.StatusInternalServerError, "Failed to delete oil record", true)
	}

	return NewHTTPResponse(c, fiber.StatusOK, map[string]any{
		"message": "Oil record deleted successfully",
	})
}

func SetupOilRouter(app *fiber.App, ctrl OilController, mw middleware.HTTPMiddleware) {
	oilGroup := app.Group(BASE_OIL_PATH, mw.Verify)

	oilGroup.Get(OIL_GET, ctrl.GetOil)
	oilGroup.Patch(OIL_UPDATE, ctrl.UpdateOil)
	oilGroup.Delete(OIL_DELETE, ctrl.DeleteOil)

	oilGroup.Get(OIL_GETMANY, ctrl.GetOilByCollectorId)
}

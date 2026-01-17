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
	BASE_REPORT_PATH = config.BASE_API_HTTP_PATH + "/reports"
	REPORT_BY_DATE   = "/date"
	REPORT_ALL       = "/all"
)

type ReportController struct {
	reportUsecase usecase.IReportUsecase
}

func NewReportController(reportUsecase usecase.IReportUsecase) ReportController {
	return ReportController{reportUsecase}
}

func (rc ReportController) GetReportByDate(c *fiber.Ctx) error {
	req := new(dto.ReportByDate)
	if err := c.BodyParser(req); err != nil {
		return fiber.ErrBadRequest
	}

	ctx := c.Context()

	result := rc.reportUsecase.GetReportByDate(ctx, req)
	if result.IsError() {
		if err := result.ExpectedError(); err != nil {
			return NewHTTPError(c, err)
		}

		log.Println(result)
		return NewHTTPErrorSimple(c, fiber.StatusInternalServerError, "Failed to get report", true)
	}

	return NewHTTPResponse(c, fiber.StatusOK, result.Value())
}

func (rc ReportController) GetAllReports(c *fiber.Ctx) error {
	req := new(dto.ReportAll)
	if err := c.BodyParser(req); err != nil {
		return fiber.ErrBadRequest
	}

	ctx := c.Context()

	result := rc.reportUsecase.GetAllReports(ctx, req)
	if result.IsError() {
		if err := result.ExpectedError(); err != nil {
			return NewHTTPError(c, err)
		}

		log.Println(result)
		return NewHTTPErrorSimple(c, fiber.StatusInternalServerError, "Failed to get all reports", true)
	}

	return NewHTTPResponse(c, fiber.StatusOK, result.Value())
}

func SetupReportRouter(app *fiber.App, ctrl ReportController, mw middleware.HTTPMiddleware) {
	app.Group(BASE_REPORT_PATH, mw.Verify).
		Post(REPORT_BY_DATE, ctrl.GetReportByDate).
		Post(REPORT_ALL, ctrl.GetAllReports)
}

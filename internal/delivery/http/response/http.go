package response

import (
	"fmt"

	. "github.com/crazydw4rf/oil-bank-backend/internal/pkg/result"
	"github.com/gofiber/fiber/v2"
)

type HTTPError struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	IsExpected bool   `json:"is_expected"`
}

type HTTPResponse[T any] struct {
	Data  T          `json:"data"`
	Error *HTTPError `json:"error,omitempty"`
}

var errorMapping = map[ErrorCause]int{
	CREDENTIALS_ERROR:      fiber.StatusUnauthorized,
	ENTITY_NOT_FOUND:       fiber.StatusNotFound,
	ENTITY_DUPLICATE:       fiber.StatusConflict,
	TOKEN_GENERATION_ERROR: fiber.StatusForbidden,
	INTERNAL_SERVICE_ERROR: fiber.StatusInternalServerError,
	BAD_REQUEST_ERROR:      fiber.StatusBadRequest,
	UNKNOWN_ERROR:          fiber.StatusInternalServerError,
}

func NewHTTPResponse[T any](ctx *fiber.Ctx, code int, data T) error {
	return ctx.Status(code).JSON(HTTPResponse[T]{
		Data:  data,
		Error: nil,
	})
}

func NewHTTPError(ctx *fiber.Ctx, err error) error {
	httpErr := buildHTTPError(err)

	fmt.Printf("Error: %#v\n", httpErr)

	return ctx.Status(httpErr.Code).JSON(HTTPResponse[any]{
		Data:  nil,
		Error: httpErr,
	})
}

func NewHTTPErrorSimple(ctx *fiber.Ctx, code int, message string, Expected ...bool) error {
	var IsExpected bool = false
	if len(Expected) > 0 {
		IsExpected = Expected[0]
	}
	return ctx.Status(code).JSON(HTTPResponse[any]{
		Data: nil,
		Error: &HTTPError{
			Code:       code,
			Message:    message,
			IsExpected: IsExpected,
		},
	})
}

func buildHTTPError(err error) *HTTPError {
	errTrace, ok := err.(*ErrorTrace)
	if !ok {
		return &HTTPError{
			Code:       fiber.StatusInternalServerError,
			Message:    "Internal Server Error",
			IsExpected: false,
		}
	}

	code, ok := errorMapping[errTrace.Cause()]
	if !ok {
		code = fiber.StatusInternalServerError
	}

	return &HTTPError{
		Code:       code,
		Message:    errTrace.Error(),
		IsExpected: errTrace.IsExpected,
	}
}

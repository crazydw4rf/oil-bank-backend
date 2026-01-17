package controller

import (
	"strconv"

	"github.com/crazydw4rf/oil-bank-backend/internal/constants"
	. "github.com/crazydw4rf/oil-bank-backend/internal/pkg/result"
	"github.com/gofiber/fiber/v2"
)

func UserIdExtractor(c *fiber.Ctx) Result[int64] {
	sub, ok := c.Locals(constants.UserIdKey).(string)
	if !ok {
		return NewError[int64]("Cannot extract user ID").WithCause(INTERNAL_LOGIC_ERROR)
	}

	id, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		return NewError[int64]("Invalid user ID").WithCause(INTERNAL_LOGIC_ERROR)
	}

	return Ok(id)
}

func CollectorIdExtractor(c *fiber.Ctx) Result[int64] {
	sub, ok := c.Locals(constants.CollectorIdKey).(string)
	if !ok {
		return NewError[int64]("Cannot extract collector ID").WithCause(INTERNAL_LOGIC_ERROR)
	}

	id, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		return NewError[int64]("Invalid collector ID").WithCause(INTERNAL_LOGIC_ERROR)
	}

	return Ok(id)
}

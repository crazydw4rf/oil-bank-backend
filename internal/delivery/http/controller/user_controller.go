package controller

import (
	"fmt"
	"log"
	"time"

	"github.com/crazydw4rf/oil-bank-backend/internal/delivery/http/middleware"
	. "github.com/crazydw4rf/oil-bank-backend/internal/delivery/http/response"
	"github.com/crazydw4rf/oil-bank-backend/internal/dto"
	"github.com/crazydw4rf/oil-bank-backend/internal/entity"
	"github.com/crazydw4rf/oil-bank-backend/internal/services/config"
	"github.com/crazydw4rf/oil-bank-backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

const (
	BASE_USER_PATH = config.BASE_API_HTTP_PATH + "/users"
	// Group /users
	USER_GETMANY  = "/"
	USER_GET      = "/:id"
	USER_UPDATE   = "/:id"
	USER_DELETE   = "/:id"
	USER_CREATE   = BASE_USER_PATH + "/"
	USER_LOGIN    = BASE_USER_PATH + "/auth/login"
	REFRESH_TOKEN = BASE_USER_PATH + "/auth/refresh"
)

type UserController struct {
	userUsecase usecase.IUserUsecase
}

func NewUserController(userUsecase usecase.IUserUsecase) UserController {
	return UserController{userUsecase}
}

func (uc UserController) UserCreate(c *fiber.Ctx) error {
	req := new(dto.UserCreateRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.ErrBadRequest
	}

	ctx := c.Context()

	result := uc.userUsecase.UserRegister(ctx, req)
	if result.IsError() {
		if err := result.ExpectedError(); err != nil {
			return NewHTTPError(c, err)
		}

		log.Println(result)
		return NewHTTPErrorSimple(c, fiber.StatusInternalServerError, "Failed to create user", true)
	}

	// FIXME: untuk membuat pengguna baru tidak perlu mengembalikan data pengguna
	return NewHTTPResponse(c, fiber.StatusCreated, result.Value())
}

func (uc UserController) UserLogin(c *fiber.Ctx) error {
	req := new(dto.UserLoginRequest)
	if err := c.BodyParser(req); err != nil {
		fmt.Println(err)
		return fiber.ErrBadRequest
	}

	ctx := c.Context()

	result := uc.userUsecase.UserLogin(ctx, req)
	if result.IsError() {
		if err := result.ExpectedError(); err != nil {
			return NewHTTPError(c, err)
		}

		// TODO: pake zap library untuk logging
		log.Println(result)

		return NewHTTPErrorSimple(c, fiber.StatusInternalServerError, "Failed to login user", true)
	}

	user := result.Value()

	setAuthCookie(user, c)

	return NewHTTPResponse(c, fiber.StatusOK, user)
}

func (uc UserController) GetUser(c *fiber.Ctx) error {
	id := UserIdExtractor(c)
	if id.IsError() {
		return NewHTTPErrorSimple(c, fiber.StatusBadRequest, "Invalid user id", true)
	}
	result := uc.userUsecase.UserFind(c.Context(), id.Value())
	if result.IsError() {
		if err := result.ExpectedError(); err != nil {
			return NewHTTPError(c, err)
		}
	}

	return NewHTTPResponse(c, fiber.StatusOK, result.Value())
}

func (uc UserController) GetUserMany(c *fiber.Ctx) error {
	panic("Not implemented")
}

func (uc UserController) UpdateUser(c *fiber.Ctx) error {
	panic("Not implemented")
}

func (uc UserController) DeleteUser(c *fiber.Ctx) error {
	panic("Not implemented")
}

func setAuthCookie(user *entity.UserWithCollector, ctx *fiber.Ctx) {
	ctx.Cookie(&fiber.Cookie{
		Name:     config.ACCESS_TOKEN_COOKIE_NAME,
		Value:    user.AccessToken,
		Path:     "/",
		Expires:  time.Now().Add(config.ACCESS_TOKEN_EXPIRATION_TIME),
		HTTPOnly: true,
		Secure:   true,
	})
	ctx.Cookie(&fiber.Cookie{
		Name:     config.REFRESH_TOKEN_COOKIE_NAME,
		Value:    user.RefreshToken,
		Path:     REFRESH_TOKEN,
		Expires:  time.Now().Add(config.REFRESH_TOKEN_EXPIRATION_TIME),
		HTTPOnly: true,
		Secure:   true,
	})
}

func SetupUserRouter(app *fiber.App, ctrl UserController, mw middleware.HTTPMiddleware) {
	app.Post(USER_CREATE, ctrl.UserCreate)
	app.Post(USER_LOGIN, ctrl.UserLogin)

	app.Group(BASE_USER_PATH, mw.Verify).
		Get(USER_GET, ctrl.GetUser).
		Get(USER_GETMANY, ctrl.GetUserMany).
		Patch(USER_UPDATE, ctrl.UpdateUser).
		Delete(USER_DELETE, ctrl.DeleteUser)

}

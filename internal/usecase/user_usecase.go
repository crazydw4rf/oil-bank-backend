package usecase

import (
	"context"
	"log"

	"github.com/crazydw4rf/oil-bank-backend/internal/auth"
	"github.com/crazydw4rf/oil-bank-backend/internal/dto"
	"github.com/crazydw4rf/oil-bank-backend/internal/entity"
	. "github.com/crazydw4rf/oil-bank-backend/internal/pkg/result"
	"github.com/crazydw4rf/oil-bank-backend/internal/repository"
	"github.com/crazydw4rf/oil-bank-backend/internal/services/config"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	UserRegister(ctx context.Context, dto *dto.UserCreateRequest) Result[*entity.User]
	UserLogin(ctx context.Context, dto *dto.UserLoginRequest) Result[*entity.UserWithCollector]
	UserUpdate(ctx context.Context, dto *dto.UserUpdateRequest) Result[*entity.User]
	UserDelete(ctx context.Context, id int64) Result[bool]
	UserFind(ctx context.Context, id int64) Result[*entity.User]
}

type UserUsecase struct {
	userRepo repository.IUserRepository
	cfg      *config.Config
}

func NewUserUsecase(userRepo repository.IUserRepository, cfg *config.Config) IUserUsecase {
	return UserUsecase{userRepo, cfg}
}

var _ IUserUsecase = (*UserUsecase)(nil)

func (uc UserUsecase) UserRegister(ctx context.Context, dto *dto.UserCreateRequest) Result[*entity.User] {
	// TODO: implementasi validasi dto dengan menggunakan library validator
	user := &entity.User{
		Username: dto.Username,
		Email:    dto.Email,
		UserType: dto.UserType,
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return NewError[*entity.User]("Password salah", true).WithCause(CREDENTIALS_ERROR)
		}

		return NewError[*entity.User]("Password hashing gagal").WithCause(UNKNOWN_ERROR)
	}
	user.PasswordHash = string(hash)

	result := uc.userRepo.Create(ctx, user)
	if result.IsError() {
		log.Println(result.Error())
		return Err(result, "Gagal membuat user baru", true)
	}

	return Ok(user)
}

func (uc UserUsecase) UserLogin(ctx context.Context, dto *dto.UserLoginRequest) Result[*entity.UserWithCollector] {
	result := uc.userRepo.FindByEmailWithCollector(ctx, dto.Email)
	if result.IsError() {
		return Err(result, "Email tidak ditemukan", true)
	}
	user := result.Value()
	if user.UserType != entity.COLLECTOR {
		return Err(result, "User selain collector tidak dapat login", true).WithCause(BAD_REQUEST_ERROR)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(dto.Password)); err != nil {
		return NewError[*entity.UserWithCollector]("Password salah", true).WithCause(CREDENTIALS_ERROR)
	}

	userWithToken, err := auth.GenerateToken(user, uc.cfg)
	if err != nil {
		return NewError[*entity.UserWithCollector]("Gagal membuat token").WithCause(TOKEN_GENERATION_ERROR)
	}

	return Ok(userWithToken)
}

func (uc UserUsecase) UserUpdate(ctx context.Context, dto *dto.UserUpdateRequest) Result[*entity.User] {
	panic("Not implemented")
}

func (uc UserUsecase) UserDelete(ctx context.Context, id int64) Result[bool] {
	panic("Not implemented")
}

func (uc UserUsecase) UserFind(ctx context.Context, id int64) Result[*entity.User] {
	result := uc.userRepo.Find(ctx, id)
	if result.IsError() {
		return Err(result, "User tidak ditemukan", true)
	}
	return Ok(result.Value())
}

// func generateToken(user *entity.User, cfg *config.Config) (*entity.User, error) {
// 	at, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
// 		Subject:   strconv.FormatInt(user.Id, 10),
// 		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
// 	}).
// 		SignedString([]byte(cfg.JWT_ACCESS_TOKEN_SECRET))
// 	if err != nil {
// 		return nil, err
// 	}

// 	rt, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
// 		Subject:   strconv.FormatInt(user.Id, 10),
// 		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
// 	}).
// 		SignedString([]byte(cfg.JWT_REFRESH_TOKEN_SECRET))
// 	if err != nil {
// 		return nil, err
// 	}

// 	user.AccessToken = at
// 	user.RefreshToken = rt

// 	return user, nil
// }

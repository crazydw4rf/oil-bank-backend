package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/crazydw4rf/oil-bank-backend/internal/entity"
	. "github.com/crazydw4rf/oil-bank-backend/internal/pkg/result"
	"github.com/crazydw4rf/oil-bank-backend/internal/services"
	"github.com/crazydw4rf/oil-bank-backend/internal/types"
	"github.com/jackc/pgx"
)

type IUserRepository interface {
	types.BaseRepository[entity.User]
	FindByEmail(ctx context.Context, email string) Result[*entity.User]
	FindByEmailWithSeller(ctx context.Context, email string) Result[*entity.UserWithSeller]
	FindByEmailWithCollector(ctx context.Context, email string) Result[*entity.UserWithCollector]
	FindByEmailWithCompany(ctx context.Context, email string) Result[*entity.UserWithCompany]
}

type UserRepository struct {
	db services.DatabaseService
}

var _ IUserRepository = (*UserRepository)(nil)

func NewUserRepository(db services.DatabaseService) IUserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) Result[*entity.User] {
	rows := r.db.QueryRowxContext(ctx, userCreate,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.UserType,
	)

	err := rows.StructScan(user)
	if err != nil {
		return handleUserError[*entity.User](err)
	}

	return Ok(user)
}

func (r *UserRepository) Find(ctx context.Context, id int64) Result[*entity.User] {
	rows := r.db.QueryRowxContext(ctx, userFind, id)
	user := new(entity.User)

	err := rows.StructScan(user)
	if err != nil {
		return handleUserError[*entity.User](err)
	}

	return Ok(user)
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) Result[*entity.User] {
	rows := r.db.QueryRowxContext(ctx, userUpdate,
		user.Id,
		user.Username,
		user.Email,
	)

	err := rows.StructScan(user)
	if err != nil {
		return handleUserError[*entity.User](err)
	}

	return Ok(user)
}

func (r *UserRepository) Delete(ctx context.Context, id int64) Result[bool] {
	res, err := r.db.ExecContext(ctx, userDelete, id)
	if err != nil {
		return handleUserError[bool](err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected <= 0 {
		return NewError[bool]("can't delete user, user not found", true).WithCause(ENTITY_NOT_FOUND)
	}

	return Ok(true)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) Result[*entity.User] {
	rows := r.db.QueryRowxContext(ctx, userFindByEmail, email)
	user := new(entity.User)

	err := rows.StructScan(user)
	if err != nil {
		return handleUserError[*entity.User](err)
	}

	return Ok(user)
}

func (r *UserRepository) FindByEmailWithSeller(ctx context.Context, email string) Result[*entity.UserWithSeller] {
	rows := r.db.QueryRowxContext(ctx, userFindByEmailWithSeller, email)
	user := new(entity.UserWithSeller)

	err := rows.StructScan(user)
	if err != nil {
		return handleUserError[*entity.UserWithSeller](err)
	}

	return Ok(user)
}

func (r *UserRepository) FindByEmailWithCollector(ctx context.Context, email string) Result[*entity.UserWithCollector] {
	rows := r.db.QueryRowxContext(ctx, userFindByEmailWithCollector, email)
	user := new(entity.UserWithCollector)

	err := rows.StructScan(user)
	if err != nil {
		return handleUserError[*entity.UserWithCollector](err)
	}

	return Ok(user)
}

func (r *UserRepository) FindByEmailWithCompany(ctx context.Context, email string) Result[*entity.UserWithCompany] {
	rows := r.db.QueryRowxContext(ctx, userFindByEmailWithCompany, email)
	user := new(entity.UserWithCompany)

	err := rows.StructScan(user)
	if err != nil {
		return handleUserError[*entity.UserWithCompany](err)
	}

	return Ok(user)
}

func handleUserError[T any](err error) Result[T] {
	var pgErr pgx.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return NewError[T]("user already exists", true).WithCause(ENTITY_DUPLICATE)
		default:
			return NewError[T]("database error: " + err.Error()).WithCause(INTERNAL_SERVICE_ERROR)
		}
	} else if errors.Is(err, sql.ErrNoRows) {
		return NewError[T]("user not found", true).WithCause(ENTITY_NOT_FOUND)
	}

	return NewError[T]("database error: " + err.Error()).WithCause(INTERNAL_SERVICE_ERROR)
}

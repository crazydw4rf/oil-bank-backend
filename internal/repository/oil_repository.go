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

type IOilRepository interface {
	types.BaseRepository[entity.Oil]
	GetByCollectorId(ctx context.Context, collectorId int64) Result[*entity.Oil]
}

type OilRepository struct {
	db services.DatabaseService
}

var _ IOilRepository = (*OilRepository)(nil)

func NewOilRepository(db services.DatabaseService) IOilRepository {
	return &OilRepository{db}
}

func (r *OilRepository) Create(ctx context.Context, oil *entity.Oil) Result[*entity.Oil] {
	rows := r.db.QueryRowxContext(ctx, oilCreate, oil.CollectorId)

	err := rows.StructScan(oil)
	if err != nil {
		return handleOilError[*entity.Oil](err)
	}

	return Ok(oil)
}

func (r *OilRepository) Find(ctx context.Context, id int64) Result[*entity.Oil] {
	rows := r.db.QueryRowxContext(ctx, oilFind, id)
	oil := new(entity.Oil)

	err := rows.StructScan(oil)
	if err != nil {
		return handleOilError[*entity.Oil](err)
	}

	return Ok(oil)
}

func (r *OilRepository) GetByCollectorId(ctx context.Context, collectorId int64) Result[*entity.Oil] {
	rows := r.db.QueryRowxContext(ctx, oilGetByCollectorId, collectorId)
	if err := rows.Err(); err != nil {
		return handleOilError[*entity.Oil](err)
	}
	oil := new(entity.Oil)
	err := rows.StructScan(oil)
	if err != nil {
		return handleOilError[*entity.Oil](err)
	}

	return Ok(oil)
}

func (r *OilRepository) Update(ctx context.Context, oil *entity.Oil) Result[*entity.Oil] {
	rows := r.db.QueryRowxContext(ctx, oilUpdate,
		oil.Id,
		oil.TotalVolume,
	)

	err := rows.StructScan(oil)
	if err != nil {
		return handleOilError[*entity.Oil](err)
	}

	return Ok(oil)
}

func (r *OilRepository) Delete(ctx context.Context, id int64) Result[bool] {
	res, err := r.db.ExecContext(ctx, oilDelete, id)
	if err != nil {
		return handleOilError[bool](err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected <= 0 {
		return NewError[bool]("oil record not found", true).WithCause(ENTITY_NOT_FOUND)
	}

	return Ok(true)
}

func handleOilError[T any](err error) Result[T] {
	var pgErr pgx.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23503":
			return NewError[T]("invalid collector_id", true).WithCause(ENTITY_NOT_FOUND)
		case "23505":
			return NewError[T]("oil record already exists for this collector", true).WithCause(ENTITY_DUPLICATE)
		default:
			return NewError[T]("database error: " + err.Error()).WithCause(INTERNAL_SERVICE_ERROR)
		}
	} else if errors.Is(err, sql.ErrNoRows) {
		return NewError[T]("oil record not found", true).WithCause(ENTITY_NOT_FOUND)
	}

	return NewError[T]("database error: " + err.Error()).WithCause(INTERNAL_SERVICE_ERROR)
}

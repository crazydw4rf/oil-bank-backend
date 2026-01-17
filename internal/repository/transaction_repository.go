package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/crazydw4rf/oil-bank-backend/internal/entity"
	. "github.com/crazydw4rf/oil-bank-backend/internal/pkg/result"
	"github.com/crazydw4rf/oil-bank-backend/internal/services"
	"github.com/jackc/pgx"
)

type ITransactionRepository interface {
	CreateSellTransaction(ctx context.Context, tx *entity.SellTransaction) Result[*entity.SellTransaction]
	CreateDistributeTransaction(ctx context.Context, tx *entity.DistributeTransaction) Result[*entity.DistributeTransaction]
	UpdateSellTransaction(ctx context.Context, id int64, volume float64, price float64) Result[*entity.SellTransaction]
	UpdateDistributeTransaction(ctx context.Context, id int64, volume float64, price float64) Result[*entity.DistributeTransaction]
	FindSellTransactionById(ctx context.Context, id int64) Result[*entity.SellTransaction]
	FindDistributeTransactionById(ctx context.Context, id int64) Result[*entity.DistributeTransaction]
	UpdateCollectorVolume(ctx context.Context, collectorId int64, volumeDelta float64) Result[bool]
}

type TransactionRepository struct {
	db services.DatabaseService
}

var _ ITransactionRepository = (*TransactionRepository)(nil)

func NewTransactionRepository(db services.DatabaseService) ITransactionRepository {
	return &TransactionRepository{db}
}

func (r TransactionRepository) CreateSellTransaction(ctx context.Context, tx *entity.SellTransaction) Result[*entity.SellTransaction] {
	rows := r.db.QueryRowxContext(ctx, sellTransactionCreate,
		tx.SellerId,
		tx.CollectorId,
		tx.Volume,
		tx.Price,
	)

	err := rows.StructScan(tx)
	if err != nil {
		return handleTransactionError[*entity.SellTransaction](err)
	}

	return Ok(tx)
}

func (r TransactionRepository) CreateDistributeTransaction(ctx context.Context, tx *entity.DistributeTransaction) Result[*entity.DistributeTransaction] {
	rows := r.db.QueryRowxContext(ctx, distributeTransactionCreate,
		tx.CollectorId,
		tx.CompanyId,
		tx.Volume,
		tx.Price,
	)

	err := rows.StructScan(tx)
	if err != nil {
		return handleTransactionError[*entity.DistributeTransaction](err)
	}

	return Ok(tx)
}

func (r TransactionRepository) UpdateSellTransaction(ctx context.Context, id int64, volume float64, price float64) Result[*entity.SellTransaction] {
	tx := &entity.SellTransaction{}
	row := r.db.QueryRowxContext(ctx, sellTransactionUpdate, id, volume, price)

	err := row.StructScan(tx)
	if err != nil {
		return handleTransactionError[*entity.SellTransaction](err)
	}

	return Ok(tx)
}

func (r TransactionRepository) UpdateDistributeTransaction(ctx context.Context, id int64, volume float64, price float64) Result[*entity.DistributeTransaction] {
	tx := &entity.DistributeTransaction{}
	row := r.db.QueryRowxContext(ctx, distributeTransactionUpdate, id, volume, price)

	err := row.StructScan(tx)
	if err != nil {
		return handleTransactionError[*entity.DistributeTransaction](err)
	}

	return Ok(tx)
}

func (r TransactionRepository) FindSellTransactionById(ctx context.Context, id int64) Result[*entity.SellTransaction] {
	tx := &entity.SellTransaction{}
	row := r.db.QueryRowxContext(ctx, sellTransactionFindById, id)

	err := row.StructScan(tx)
	if err != nil {
		return handleTransactionError[*entity.SellTransaction](err)
	}

	return Ok(tx)
}

func (r TransactionRepository) FindDistributeTransactionById(ctx context.Context, id int64) Result[*entity.DistributeTransaction] {
	tx := &entity.DistributeTransaction{}
	row := r.db.QueryRowxContext(ctx, distributeTransactionFindById, id)

	err := row.StructScan(tx)
	if err != nil {
		return handleTransactionError[*entity.DistributeTransaction](err)
	}

	return Ok(tx)
}

func (r TransactionRepository) UpdateCollectorVolume(ctx context.Context, collectorId int64, volumeDelta float64) Result[bool] {
	_, err := r.db.ExecContext(ctx, updateCollectorVolume, collectorId, volumeDelta)
	if err != nil {
		return handleTransactionError[bool](err)
	}

	return Ok(true)
}

func handleTransactionError[T any](err error) Result[T] {
	var pgErr pgx.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23503":
			return NewError[T]("referenced entity not found", true).WithCause(ENTITY_NOT_FOUND)
		case "23514":
			return NewError[T]("constraint violation", true).WithCause(INTERNAL_SERVICE_ERROR)
		default:
			return NewError[T]("database error: " + err.Error()).WithCause(INTERNAL_SERVICE_ERROR)
		}
	} else if errors.Is(err, sql.ErrNoRows) {
		return NewError[T]("transaction not found", true).WithCause(ENTITY_NOT_FOUND)
	}

	return NewError[T]("database error: " + err.Error()).WithCause(INTERNAL_SERVICE_ERROR)
}

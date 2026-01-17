package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/crazydw4rf/oil-bank-backend/internal/dto"
	"github.com/crazydw4rf/oil-bank-backend/internal/entity"
	. "github.com/crazydw4rf/oil-bank-backend/internal/pkg/result"
	"github.com/crazydw4rf/oil-bank-backend/internal/repository"
)

type ITransactionUsecase interface {
	CreateTransaction(ctx context.Context, collectorId int64, dto *dto.TransactionCreateDto) Result[*dto.TransactionResponse]
	UpdateTransaction(ctx context.Context, collectorId int64, id int64, updateDto *dto.UpdateTransactionDto) Result[*dto.TransactionResponse]
}

type TransactionUsecase struct {
	transactionRepo repository.ITransactionRepository
	userRepo        repository.IUserRepository
}

func NewTransactionUsecase(transactionRepo repository.ITransactionRepository, userRepo repository.IUserRepository) ITransactionUsecase {
	return &TransactionUsecase{transactionRepo, userRepo}
}

var _ ITransactionUsecase = (*TransactionUsecase)(nil)

func (uc *TransactionUsecase) CreateTransaction(ctx context.Context, collectorId int64, dtoo *dto.TransactionCreateDto) Result[*dto.TransactionResponse] {
	switch dtoo.TransactionType {
	case dto.TRANSACTION_SELL:
		return uc.createSellTransaction(ctx, collectorId, dtoo)
	case dto.TRANSACTION_BUY:
		return uc.createDistributeTransaction(ctx, collectorId, dtoo)
	default:
		return NewError[*dto.TransactionResponse]("Invalid transaction type", true).WithCause(BAD_REQUEST_ERROR)
	}
}

func (uc *TransactionUsecase) createSellTransaction(ctx context.Context, collectorId int64, txDto *dto.TransactionCreateDto) Result[*dto.TransactionResponse] {
	result := uc.userRepo.FindByEmailWithSeller(ctx, txDto.Email)
	if result.IsError() {
		return NewError[*dto.TransactionResponse]("User seller not found", true).WithCause(ENTITY_NOT_FOUND)
	}
	userWithSeller := result.Value()

	if txDto.OilVolume <= 0 {
		return NewError[*dto.TransactionResponse]("Volume must be greater than 0", true).WithCause(INTERNAL_LOGIC_ERROR)
	}
	if txDto.Price <= 0 {
		return NewError[*dto.TransactionResponse]("Price must be greater than 0", true).WithCause(INTERNAL_LOGIC_ERROR)
	}

	tx := &entity.SellTransaction{
		SellerId:    userWithSeller.SellerId,
		CollectorId: collectorId,
		Price:       txDto.Price,
		Volume:      txDto.OilVolume,
	}

	res := uc.transactionRepo.CreateSellTransaction(ctx, tx)
	if res.IsError() {
		log.Println(res.Error())
		return NewError[*dto.TransactionResponse](
			"Failed to create sell transaction",
			true,
		).WithCause(INTERNAL_SERVICE_ERROR)
	}

	transaction := res.Value()
	response := &dto.TransactionResponse{
		Id:              transaction.Id,
		SellerId:        transaction.SellerId,
		OilVolume:       transaction.Volume,
		Price:           transaction.Price,
		TransactionType: dto.TRANSACTION_SELL,
		CreatedAt:       transaction.CreatedAt,
		UpdatedAt:       transaction.UpdatedAt,
	}

	return Ok(response)
}

func (uc *TransactionUsecase) createDistributeTransaction(ctx context.Context, collectorId int64, txDto *dto.TransactionCreateDto) Result[*dto.TransactionResponse] {
	result := uc.userRepo.FindByEmailWithCompany(ctx, txDto.Email)
	if e := result.LastError(); e != nil {
		return NewError[*dto.TransactionResponse](e.Error(), e.IsExpected).WithCause(e.Cause())
	}
	userWithCompany := result.Value()

	if txDto.OilVolume <= 0 {
		return NewError[*dto.TransactionResponse]("Volume must be greater than 0", true).WithCause(INTERNAL_LOGIC_ERROR)
	}
	if txDto.Price <= 0 {
		return NewError[*dto.TransactionResponse]("Price must be greater than 0", true).WithCause(INTERNAL_LOGIC_ERROR)
	}

	tx := &entity.DistributeTransaction{
		CollectorId: collectorId,
		CompanyId:   userWithCompany.CompanyId,
		Volume:      txDto.OilVolume,
		Price:       txDto.Price,
	}

	res := uc.transactionRepo.CreateDistributeTransaction(ctx, tx)
	if res.IsError() {
		log.Println(res.Error())

		if res.RootError() != nil && res.RootError().Cause() == INTERNAL_SERVICE_ERROR {
			return NewError[*dto.TransactionResponse](
				"Failed to create distribute transaction. Insufficient oil inventory or invalid company ID",
				true,
			).WithCause(INTERNAL_SERVICE_ERROR)
		}

		return NewError[*dto.TransactionResponse](
			"Failed to create distribute transaction",
			true,
		).WithCause(INTERNAL_SERVICE_ERROR)
	}

	transaction := res.Value()
	response := &dto.TransactionResponse{
		Id:              transaction.Id,
		CompanyId:       transaction.CompanyId,
		OilVolume:       transaction.Volume,
		Price:           transaction.Price,
		TransactionType: dto.TRANSACTION_BUY,
		CreatedAt:       transaction.CreatedAt,
		UpdatedAt:       transaction.UpdatedAt,
	}

	return Ok(response)
}

func (uc *TransactionUsecase) UpdateTransaction(ctx context.Context, collectorId int64, id int64, updateDto *dto.UpdateTransactionDto) Result[*dto.TransactionResponse] {
	if updateDto.OilVolume <= 0 {
		return NewError[*dto.TransactionResponse]("Volume must be greater than 0", true).WithCause(INTERNAL_LOGIC_ERROR)
	}
	if updateDto.Price <= 0 {
		return NewError[*dto.TransactionResponse]("Price must be greater than 0", true).WithCause(INTERNAL_LOGIC_ERROR)
	}

	switch updateDto.TransactionType {
	case dto.TRANSACTION_SELL:
		return uc.updateSellTransaction(ctx, collectorId, id, updateDto)
	case dto.TRANSACTION_BUY:
		return uc.updateDistributeTransaction(ctx, collectorId, id, updateDto)
	default:
		return NewError[*dto.TransactionResponse]("Invalid transaction type", true).WithCause(BAD_REQUEST_ERROR)
	}
}

func (uc *TransactionUsecase) updateSellTransaction(ctx context.Context, collectorId int64, id int64, updateDto *dto.UpdateTransactionDto) Result[*dto.TransactionResponse] {
	// First, verify the transaction exists
	findResult := uc.transactionRepo.FindSellTransactionById(ctx, id)
	if findResult.IsError() {
		return NewError[*dto.TransactionResponse](
			fmt.Sprintf("Sell transaction with id %d not found", id),
			true,
		).WithCause(ENTITY_NOT_FOUND)
	}

	// Verify the transaction belongs to this collector
	existingTransaction := findResult.Value()
	if existingTransaction.CollectorId != collectorId {
		return NewError[*dto.TransactionResponse](
			"You are not authorized to update this transaction",
			true,
		).WithCause(CREDENTIALS_ERROR)
	}

	// Update the transaction
	result := uc.transactionRepo.UpdateSellTransaction(ctx, id, updateDto.OilVolume, updateDto.Price)
	if result.IsError() {
		log.Println(result.Error())
		return NewError[*dto.TransactionResponse](
			"Failed to update sell transaction",
			true,
		).WithCause(INTERNAL_SERVICE_ERROR)
	}

	transaction := result.Value()
	response := &dto.TransactionResponse{
		Id:              transaction.Id,
		SellerId:        transaction.SellerId,
		OilVolume:       transaction.Volume,
		Price:           transaction.Price,
		TransactionType: dto.TRANSACTION_SELL,
		CreatedAt:       transaction.CreatedAt,
		UpdatedAt:       transaction.UpdatedAt,
	}

	return Ok(response)
}

func (uc *TransactionUsecase) updateDistributeTransaction(ctx context.Context, collectorId int64, id int64, updateDto *dto.UpdateTransactionDto) Result[*dto.TransactionResponse] {
	// First, verify the transaction exists
	findResult := uc.transactionRepo.FindDistributeTransactionById(ctx, id)
	if findResult.IsError() {
		return NewError[*dto.TransactionResponse](
			fmt.Sprintf("Distribute transaction with id %d not found", id),
			true,
		).WithCause(ENTITY_NOT_FOUND)
	}

	// Verify the transaction belongs to this collector
	existingTransaction := findResult.Value()
	if existingTransaction.CollectorId != collectorId {
		return NewError[*dto.TransactionResponse](
			"You are not authorized to update this transaction",
			true,
		).WithCause(CREDENTIALS_ERROR)
	}

	// Update the transaction
	result := uc.transactionRepo.UpdateDistributeTransaction(ctx, id, updateDto.OilVolume, updateDto.Price)
	if result.IsError() {
		log.Println(result.Error())
		return NewError[*dto.TransactionResponse](
			"Failed to update distribute transaction",
			true,
		).WithCause(INTERNAL_SERVICE_ERROR)
	}

	transaction := result.Value()
	response := &dto.TransactionResponse{
		Id:              transaction.Id,
		CompanyId:       transaction.CompanyId,
		OilVolume:       transaction.Volume,
		Price:           transaction.Price,
		TransactionType: dto.TRANSACTION_BUY,
		CreatedAt:       transaction.CreatedAt,
		UpdatedAt:       transaction.UpdatedAt,
	}

	return Ok(response)
}

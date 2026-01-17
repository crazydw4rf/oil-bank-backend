package usecase

import (
	"context"

	"github.com/crazydw4rf/oil-bank-backend/internal/dto"
	"github.com/crazydw4rf/oil-bank-backend/internal/entity"
	. "github.com/crazydw4rf/oil-bank-backend/internal/pkg/result"
	"github.com/crazydw4rf/oil-bank-backend/internal/repository"
)

type IOilUsecase interface {
	GetOil(ctx context.Context, id int64) Result[*dto.OilResponse]
	GetOilByCollectorId(ctx context.Context, collectorId int64) Result[*dto.OilResponse]
	UpdateOil(ctx context.Context, id int64, totalVolume float64) Result[*dto.OilResponse]
	DeleteOil(ctx context.Context, id int64) Result[bool]
}

type OilUsecase struct {
	oilRepo repository.IOilRepository
}

func NewOilUsecase(oilRepo repository.IOilRepository) IOilUsecase {
	return &OilUsecase{oilRepo}
}

var _ IOilUsecase = (*OilUsecase)(nil)

func (uc *OilUsecase) GetOil(ctx context.Context, id int64) Result[*dto.OilResponse] {
	result := uc.oilRepo.Find(ctx, id)
	if result.IsError() {
		return NewError[*dto.OilResponse]("Failed to get oil record", true).WithCause(result.RootError().Cause())
	}

	return Ok(mapOilToResponse(result.Value()))
}

func (uc *OilUsecase) GetOilByCollectorId(ctx context.Context, collectorId int64) Result[*dto.OilResponse] {
	result := uc.oilRepo.GetByCollectorId(ctx, collectorId)
	if result.IsError() {
		return NewError[*dto.OilResponse]("Failed to get oil inventory for collector", true).WithCause(result.RootError().Cause())
	}

	return Ok(mapOilToResponse(result.Value()))
}

func (uc *OilUsecase) UpdateOil(ctx context.Context, id int64, totalVolume float64) Result[*dto.OilResponse] {

	if totalVolume < 0 {
		return NewError[*dto.OilResponse]("Total volume cannot be negative", true).WithCause(INTERNAL_LOGIC_ERROR)
	}

	existingResult := uc.oilRepo.Find(ctx, id)
	if existingResult.IsError() {
		return NewError[*dto.OilResponse]("Failed to find oil record", true).WithCause(existingResult.RootError().Cause())
	}

	oil := existingResult.Value()
	oil.TotalVolume = totalVolume

	// if !oil.IsValid() {
	// 	return NewError[*dto.OilResponse]("Invalid oil data", true).WithCause(INTERNAL_LOGIC_ERROR)
	// }

	result := uc.oilRepo.Update(ctx, oil)
	if result.IsError() {
		return NewError[*dto.OilResponse]("Failed to update oil record", true).WithCause(result.RootError().Cause())
	}

	return Ok(mapOilToResponse(result.Value()))
}

func (uc *OilUsecase) DeleteOil(ctx context.Context, id int64) Result[bool] {
	result := uc.oilRepo.Delete(ctx, id)
	if result.IsError() {
		return Err(result, "Failed to delete oil record", true)
	}

	return Ok(true)
}

func mapOilToResponse(oil *entity.Oil) *dto.OilResponse {
	return &dto.OilResponse{
		Id:          oil.Id,
		TotalVolume: oil.TotalVolume,
		CreatedAt:   oil.CreatedAt,
		UpdatedAt:   oil.UpdatedAt,
	}
}

// func mapOilToResponseMany(oils []entity.Oil) []dto.OilResponse {
// 	var responses []dto.OilResponse
// 	for _, oil := range oils {
// 		responses = append(responses, *mapOilToResponse(&oil))
// 	}
// 	return responses
// }

package usecase

import (
	"context"
	"log"

	"github.com/crazydw4rf/oil-bank-backend/internal/dto"
	"github.com/crazydw4rf/oil-bank-backend/internal/entity"
	. "github.com/crazydw4rf/oil-bank-backend/internal/pkg/result"
	"github.com/crazydw4rf/oil-bank-backend/internal/repository"
)

type IReportUsecase interface {
	GetReportByDate(ctx context.Context, dto *dto.ReportByDate) Result[[]entity.ReportTransaction]
	GetAllReports(ctx context.Context, dto *dto.ReportAll) Result[[]entity.ReportTransaction]
}

type ReportUsecase struct {
	reportRepo repository.IReportRepository
}

func NewReportUsecase(reportRepo repository.IReportRepository) IReportUsecase {
	return &ReportUsecase{reportRepo}
}

var _ IReportUsecase = (*ReportUsecase)(nil)

func (uc *ReportUsecase) GetReportByDate(ctx context.Context, reportDto *dto.ReportByDate) Result[[]entity.ReportTransaction] {
	startDate := reportDto.StartDate.Format("2006-01-02")
	endDate := reportDto.EndDate.Format("2006-01-02")

	switch reportDto.ReportType {
	case dto.REPORT_SALES:
		result := uc.reportRepo.GetSalesReport(ctx, startDate, endDate)
		if result.IsError() {
			log.Println(result.Error())
			return Err(result, "Failed to get sales report", true)
		}
		return result
	case dto.REPORT_PURCHASE:
		result := uc.reportRepo.GetPurchasesReport(ctx, startDate, endDate)
		if result.IsError() {
			log.Println(result.Error())
			return Err(result, "Failed to get purchases report", true)
		}
		return result
	}

	return NewError[[]entity.ReportTransaction]("Invalid report type", true).WithCause(UNKNOWN_ERROR)
}

func (uc *ReportUsecase) GetAllReports(ctx context.Context, reportDto *dto.ReportAll) Result[[]entity.ReportTransaction] {
	switch reportDto.ReportType {
	case dto.REPORT_SALES:
		result := uc.reportRepo.GetAllSalesReport(ctx)
		if result.IsError() {
			log.Println(result.Error())
			return Err(result, "Failed to get all sales report", true)
		}
		return result
	case dto.REPORT_PURCHASE:
		result := uc.reportRepo.GetAllPurchasesReport(ctx)
		if result.IsError() {
			log.Println(result.Error())
			return Err(result, "Failed to get all purchases report", true)
		}
		return result
	}

	return NewError[[]entity.ReportTransaction]("Invalid report type", true).WithCause(UNKNOWN_ERROR)
}

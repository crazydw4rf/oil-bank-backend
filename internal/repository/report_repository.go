package repository

import (
	"context"

	"github.com/crazydw4rf/oil-bank-backend/internal/entity"
	. "github.com/crazydw4rf/oil-bank-backend/internal/pkg/result"
	"github.com/crazydw4rf/oil-bank-backend/internal/services"
)

type IReportRepository interface {
	GetSalesReport(ctx context.Context, startDate, endDate string) Result[[]entity.ReportTransaction]
	GetPurchasesReport(ctx context.Context, startDate, endDate string) Result[[]entity.ReportTransaction]
	GetAllSalesReport(ctx context.Context) Result[[]entity.ReportTransaction]
	GetAllPurchasesReport(ctx context.Context) Result[[]entity.ReportTransaction]
}

type ReportRepository struct {
	db services.DatabaseService
}

var _ IReportRepository = (*ReportRepository)(nil)

func NewReportRepository(db services.DatabaseService) IReportRepository {
	return &ReportRepository{db}
}

func (r ReportRepository) GetSalesReport(ctx context.Context, startDate, endDate string) Result[[]entity.ReportTransaction] {
	rows, err := r.db.QueryxContext(ctx, reportSalesByDate, startDate, endDate)
	if err != nil {
		return handleTransactionError[[]entity.ReportTransaction](err)
	}
	defer rows.Close()

	var reports []entity.ReportTransaction
	for rows.Next() {
		var report entity.ReportTransaction
		err := rows.StructScan(&report)
		if err != nil {
			return handleTransactionError[[]entity.ReportTransaction](err)
		}
		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return handleTransactionError[[]entity.ReportTransaction](err)
	}

	return Ok(reports)
}

func (r ReportRepository) GetPurchasesReport(ctx context.Context, startDate, endDate string) Result[[]entity.ReportTransaction] {
	rows, err := r.db.QueryxContext(ctx, reportPurchasesByDate, startDate, endDate)
	if err != nil {
		return handleTransactionError[[]entity.ReportTransaction](err)
	}
	defer rows.Close()

	var reports []entity.ReportTransaction
	for rows.Next() {
		var report entity.ReportTransaction
		err := rows.StructScan(&report)
		if err != nil {
			return handleTransactionError[[]entity.ReportTransaction](err)
		}
		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return handleTransactionError[[]entity.ReportTransaction](err)
	}

	return Ok(reports)
}

func (r ReportRepository) GetAllSalesReport(ctx context.Context) Result[[]entity.ReportTransaction] {
	rows, err := r.db.QueryxContext(ctx, reportAllSales)
	if err != nil {
		return handleTransactionError[[]entity.ReportTransaction](err)
	}
	defer rows.Close()

	var reports []entity.ReportTransaction
	for rows.Next() {
		var report entity.ReportTransaction
		err := rows.StructScan(&report)
		if err != nil {
			return handleTransactionError[[]entity.ReportTransaction](err)
		}
		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return handleTransactionError[[]entity.ReportTransaction](err)
	}

	return Ok(reports)
}

func (r ReportRepository) GetAllPurchasesReport(ctx context.Context) Result[[]entity.ReportTransaction] {
	rows, err := r.db.QueryxContext(ctx, reportAllPurchases)
	if err != nil {
		return handleTransactionError[[]entity.ReportTransaction](err)
	}
	defer rows.Close()

	var reports []entity.ReportTransaction
	for rows.Next() {
		var report entity.ReportTransaction
		err := rows.StructScan(&report)
		if err != nil {
			return handleTransactionError[[]entity.ReportTransaction](err)
		}
		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return handleTransactionError[[]entity.ReportTransaction](err)
	}

	return Ok(reports)
}

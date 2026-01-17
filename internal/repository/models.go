package repository

import (
	"time"
)

type UserModel struct {
	Id           int64     `db:"id"`
	AddressId    *int64    `db:"address_id"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	UserType     string    `db:"user_type"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type AddressModel struct {
	Id            int64     `db:"id"`
	StreetAddress string    `db:"street_address"`
	City          string    `db:"city"`
	Regency       *string   `db:"regency"`
	Province      *string   `db:"province"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type SellerModel struct {
	Id         int64  `db:"id"`
	UserId     int64  `db:"user_id"`
	SellerName string `db:"seller_name"`
}

type CollectorModel struct {
	Id            int64     `db:"id"`
	UserId        int64     `db:"user_id"`
	CollectorName string    `db:"collector_name"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type CompanyModel struct {
	Id          int64  `db:"id"`
	UserId      int64  `db:"user_id"`
	CompanyName string `db:"company_name"`
}

type OilModel struct {
	Id          int64     `db:"id"`
	CollectorId int64     `db:"collector_id"`
	TotalVolume float64   `db:"total_volume"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type SellTransactionModel struct {
	Id          int64     `db:"id"`
	SellerId    int64     `db:"seller_id"`
	CollectorId int64     `db:"collector_id"`
	Volume      float64   `db:"volume"`
	Price       float64   `db:"price"`
	CreatedAt   time.Time `db:"created_at"`
}

type DistributeTransactionModel struct {
	Id          int64     `db:"id"`
	CollectorId int64     `db:"collector_id"`
	CompanyId   int64     `db:"company_id"`
	Volume      float64   `db:"volume"`
	Price       float64   `db:"price"`
	CreatedAt   time.Time `db:"created_at"`
}

type CollectorInventorySummary struct {
	CollectorId int64   `db:"collector_id"`
	TotalVolume float64 `db:"total_volume"`
}

type TransactionHistoryView struct {
	Id              int64     `db:"id"`
	TransactionType string    `db:"transaction_type"`
	SellerId        *int64    `db:"seller_id"`
	CollectorId     int64     `db:"collector_id"`
	CompanyId       *int64    `db:"company_id"`
	Volume          float64   `db:"volume"`
	Price           float64   `db:"price"`
	TotalAmount     float64   `db:"total_amount"`
	CreatedAt       time.Time `db:"created_at"`
}

type ReportSummary struct {
	TotalTransactions int       `db:"total_transactions"`
	TotalVolume       float64   `db:"total_volume"`
	TotalAmount       float64   `db:"total_amount"`
	AveragePrice      float64   `db:"average_price"`
	MinPrice          float64   `db:"min_price"`
	MaxPrice          float64   `db:"max_price"`
	StartDate         time.Time `db:"start_date"`
	EndDate           time.Time `db:"end_date"`
}

type OilInventoryDetail struct {
	Id            int64     `db:"id"`
	CollectorId   int64     `db:"collector_id"`
	CollectorName string    `db:"collector_name"`
	TotalVolume   float64   `db:"total_volume"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

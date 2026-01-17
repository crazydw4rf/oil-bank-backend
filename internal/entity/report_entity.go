package entity

import "time"

type ReportTransaction struct {
	TransactionDate time.Time `db:"transaction_date" json:"transaction_date"`
	CollectorName   string    `db:"collector_name" json:"collector_name"`
	SellerName      string    `db:"seller_name" json:"seller_name,omitempty"`
	CompanyName     string    `db:"company_name" json:"company_name,omitempty"`
	OilVolume       float64   `db:"oil_volume" json:"oil_volume"`
	Price           float64   `db:"price" json:"price"`
}

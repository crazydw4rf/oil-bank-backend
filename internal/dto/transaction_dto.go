package dto

import (
	"time"
)

type TransactionType string

const (
	TRANSACTION_SELL TransactionType = "SELL"
	TRANSACTION_BUY  TransactionType = "BUY"
)

type TransactionCreateDto struct {
	Email           string          `json:"email"`
	OilVolume       float64         `json:"oil_volume"`
	Price           float64         `json:"price"`
	TransactionType TransactionType `json:"transaction_type"`
}

type UpdateTransactionDto struct {
	TransactionType TransactionType `json:"transaction_type" validate:"required"`
	OilVolume       float64         `json:"oil_volume" validate:"required,gt=0"`
	Price           float64         `json:"price" validate:"required,gt=0"`
}

type TransactionResponse struct {
	Id              int64           `json:"id"`
	SellerId        int64           `json:"seller_id,omitempty"`
	CompanyId       int64           `json:"company_id,omitempty"`
	OilVolume       float64         `json:"oil_volume"`
	Price           float64         `json:"price"`
	TransactionType TransactionType `json:"transaction_type"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// - TextField email
// - TextField volume minyak (liter)
// - TextField Harga
// - ComboBox (Jual, Beli)
// - Button (tergantung pilihan combo box Jual/Beli)

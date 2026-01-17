package entity

import (
	"time"
)

type SellTransaction struct {
	Id          int64     `db:"id" json:"id"`
	SellerId    int64     `db:"seller_id" json:"seller_id" validate:"required"`
	CollectorId int64     `db:"collector_id" json:"collector_id" validate:"required"`
	Volume      float64   `db:"volume" json:"volume" validate:"required,gt=0"`
	Price       float64   `db:"price" json:"price" validate:"required,gt=0"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

func (st *SellTransaction) IsValid() bool {
	return st.SellerId > 0 &&
		st.CollectorId > 0 &&
		st.Volume > 0 &&
		st.Price > 0
}

func (st *SellTransaction) CalculateTotalAmount() float64 {
	return st.Volume * st.Price
}

func (st *SellTransaction) GetTotalAmount() float64 {
	return st.CalculateTotalAmount()
}

type DistributeTransaction struct {
	Id          int64     `db:"id" json:"id"`
	CollectorId int64     `db:"collector_id" json:"collector_id" validate:"required"`
	CompanyId   int64     `db:"company_id" json:"company_id" validate:"required"`
	Volume      float64   `db:"volume" json:"volume" validate:"required,gt=0"`
	Price       float64   `db:"price" json:"price" validate:"required,gt=0"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

func (dt *DistributeTransaction) IsValid() bool {
	return dt.CollectorId > 0 &&
		dt.CompanyId > 0 &&
		dt.Volume > 0 &&
		dt.Price > 0
}

func (dt *DistributeTransaction) CalculateTotalAmount() float64 {
	return dt.Volume * dt.Price
}

func (dt *DistributeTransaction) GetTotalAmount() float64 {
	return dt.CalculateTotalAmount()
}

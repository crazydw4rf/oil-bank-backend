package entity

import (
	"time"
)

type Oil struct {
	Id          int64     `db:"id" json:"id"`
	CollectorId int64     `db:"collector_id" json:"collector_id" validate:"required"`
	TotalVolume float64   `db:"total_volume" json:"total_volume" validate:"gte=0"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// func (o *Oil) IsValid() bool {
// 	return o.CollectorId > 0 && o.TotalVolume >= 0
// }

func (o *Oil) HasSufficientVolume(amount float64) bool {
	return o.TotalVolume >= amount
}

func (o *Oil) CanReduce(amount float64) bool {
	return amount > 0 && amount <= o.TotalVolume
}

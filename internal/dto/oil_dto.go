package dto

import "time"

type OilResponse struct {
	Id          int64     `json:"id"`
	TotalVolume float64   `json:"total_volume"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

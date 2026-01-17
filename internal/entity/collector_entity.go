package entity

import (
	"time"
)

type Collector struct {
	Id            int64     `db:"id" json:"id"`
	UserId        int64     `db:"user_id" json:"user_id"`
	CollectorName string    `db:"collector_name" json:"collector_name" validate:"required"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

func (c *Collector) IsValid() bool {
	return c.UserId > 0 && c.CollectorName != ""
}

package entity

import (
	"time"
)

type UserType string

const (
	SELLER    UserType = "SELLER"
	COLLECTOR UserType = "COLLECTOR"
	COMPANY   UserType = "COMPANY"
)

type User struct {
	Id           int64       `db:"id" json:"id"`
	AddressId    *int64      `db:"address_id,omitempty" json:"-"`
	Username     string      `db:"username" json:"username" validate:"required,min=7"`
	Email        string      `db:"email" json:"email,omitempty" validate:"required,email"`
	PasswordHash string      `db:"password_hash" json:"-"`
	Address      UserAddress `db:"address" json:"address"`
	UserType     UserType    `db:"user_type" json:"user_type" validate:"required"`
	CreatedAt    time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time   `db:"updated_at" json:"updated_at"`
	AccessToken  string      `db:"-" json:"access_token,omitempty"`
	RefreshToken string      `db:"-" json:"-"`
}

type UserWithSeller struct {
	User
	SellerName string `json:"seller_name" db:"seller_name"`
	SellerId   int64  `json:"seller_id" db:"seller_id"`
}

type UserWithCollector struct {
	User
	CollectorName string `json:"collector_name" db:"collector_name"`
	CollectorId   int64  `json:"collector_id" db:"collector_id"`
}

type UserWithCompany struct {
	User
	CompanyName string `json:"company_name" db:"company_name"`
	CompanyId   int64  `json:"company_id" db:"company_id"`
}

type UserAddress struct {
	Id            int64     `db:"id" json:"id"`
	StreetAddress string    `db:"street_address" json:"street_address"`
	City          string    `db:"city" json:"city"`
	Regency       string    `db:"regency" json:"regency"`
	Province      string    `db:"province" json:"province"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

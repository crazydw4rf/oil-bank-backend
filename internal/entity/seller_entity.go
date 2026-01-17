package entity

type Seller struct {
	Id         int64  `db:"id" json:"id"`
	UserId     int64  `db:"user_id" json:"user_id"`
	SellerName string `db:"seller_name" json:"seller_name" validate:"required"`
}

func (s *Seller) IsValid() bool {
	return s.UserId > 0 && s.SellerName != ""
}

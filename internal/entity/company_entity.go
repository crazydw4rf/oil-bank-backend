package entity

type Company struct {
	Id          int64  `db:"id" json:"id"`
	UserId      int64  `db:"user_id" json:"user_id"`
	CompanyName string `db:"company_name" json:"company_name" validate:"required"`
}

func (c *Company) IsValid() bool {
	return c.UserId > 0 && c.CompanyName != ""
}

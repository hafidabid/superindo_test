package model

import "time"

type Product struct {
	ID           int       `gorm:"primaryKey" json:"id"`
	ProductName  string    `gorm:"size:255" json:"product_name"`
	Qty          uint      `gorm:"type:numeric;default:0" json:"qty"`
	SellingPrice float64   `gorm:"type:numeric(10,2);default:0" json:"selling_price"`
	PromoPrice   float64   `gorm:"type:numeric(10,2);default:0" json:"promo_price"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:null" json:"updated_at"`
}

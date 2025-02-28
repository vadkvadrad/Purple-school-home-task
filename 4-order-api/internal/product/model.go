package product

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

const (
	CurrencyRUB = "RUB"
)

type Product struct {
	gorm.Model
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Images      pq.StringArray `json:"images" gorm:"type:text"`
	Price       int            `json:"price"`
	Currency    string         `json:"currency"`
	Owner       string         `json:"owner"`
}

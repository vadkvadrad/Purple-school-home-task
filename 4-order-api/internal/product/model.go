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
	Carts       pq.Int64Array  `json:"carts" gorm:"type:text"`
}

func (prod *Product) CartRemove(idToRemove int64) {
	newArray := pq.Int64Array{}
	for _, val := range prod.Carts {
		if val != idToRemove {
			newArray = append(newArray, val)
		}
	}
	prod.Carts = newArray
}

func (prod *Product) CartAdd(idToAdd int64) {
	carts := prod.Carts
	carts = append(carts, idToAdd)
	prod.Carts = carts
}


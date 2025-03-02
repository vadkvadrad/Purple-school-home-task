package cart

import (
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	Phone    string         `json:"phone"`
	Products pq.Int64Array  `json:"products" gorm:"type:text"`
	Date     datatypes.Date `json:"date"`
}

type ICartRepository interface {
	Create(cart *Cart) (*Cart, error)
	Update(cart *Cart) (*Cart, error)
	Delete(id uint64) error
	FindByID(id uint64) (*Cart, error)
	FindByPhone(phone string) []Cart
}

type ICartService interface {
	Create(cart *Cart) (*Cart, error)
	GetByID(id uint64, phone string) (*Cart, error)
	GetByPhone(phone string) []Cart
	Update(id uint64, newCart *Cart) (*Cart, error)
	Delete(id uint64, phone string) (*Cart, error)
}

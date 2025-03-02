package product

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

const (
	CurrencyRUB = "RUB"
)

type IProductRepository interface {
	Create(prod *Product) (*Product, error)
	Update(prod *Product) (*Product, error)
	Delete(id uint64) error
	FindById(id uint64) (*Product, error)
	FindByIds(ids pq.Int64Array) ([]Product, error)
	FindByIdUnscoped(id uint64) (*Product, error)
	AddMark(products []Product, idToAdd uint64) error
	DeleteMark(products []Product, idToDelete uint64) error
}

type IProductService interface {
	Create(prod *Product) (*Product, error)
	GetByID(id uint64) (*Product, error)
	Update(phone string, prod *Product) (*Product, error)
	Delete(owner, user string, id uint64) error
	GetByIDs(cart pq.Int64Array) ([]Product, error)
}


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
	for _, val := range prod.Carts {
		if val == idToAdd {
			return
		}
	}
	carts = append(carts, idToAdd)
	prod.Carts = carts
}


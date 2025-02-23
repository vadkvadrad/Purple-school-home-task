package cart

import (
	"order-api/pkg/db"

	"gorm.io/gorm/clause"
)

type CartRepository struct {
	Database *db.Db
}

func NewCartRepository(db *db.Db) *CartRepository {
	return &CartRepository{
		Database: db,
	}
}

func (repo *CartRepository) Create(cart *Cart) (*Cart, error) {
	result := repo.Database.Create(cart)
	if result.Error != nil {
		return nil, result.Error
	}
	return cart, nil
}

func (repo *CartRepository) Update(cart *Cart) (*Cart, error) {
	result := repo.Database.Clauses(clause.Returning{}).Updates(cart)
	if result.Error != nil {
		return nil, result.Error
	}
	return cart, nil
} 

func (repo *CartRepository) FindByPhone(phone string) (*Cart, error) {
	cart := &Cart{}
	result := repo.Database.First(cart, "phone = ?", phone)
	if result.Error != nil {
		return nil, result.Error
	}
	return cart, nil
}

func (repo *CartRepository) GetAll(limit, offset int) []Cart {
	var carts []Cart

	repo.Database.
		Table("links").
		Where("deleted_at is null").
		Order("id asc").
		Limit(limit).
		Offset(offset). 
		Scan(&carts)
	return carts
}
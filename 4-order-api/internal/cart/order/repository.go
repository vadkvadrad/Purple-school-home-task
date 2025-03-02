package order

import (
	"order-api/internal/cart"
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

func (repo *CartRepository) Create(cart *cart.Cart) (*cart.Cart, error) {
	result := repo.Database.Create(cart)
	if result.Error != nil {
		return nil, result.Error
	}
	return cart, nil
}

func (repo *CartRepository) Update(cart *cart.Cart) (*cart.Cart, error) {
	result := repo.Database.Clauses(clause.Returning{}).Updates(cart)
	if result.Error != nil {
		return nil, result.Error
	}
	return cart, nil
}

func (repo *CartRepository) Delete(id uint64) error {
	result := repo.Database.Delete(&cart.Cart{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *CartRepository) FindByID(id uint64) (*cart.Cart, error) {
	cart := &cart.Cart{}
	result := repo.Database.First(cart, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return cart, nil
}

func (repo *CartRepository) FindByPhone(phone string) []cart.Cart {
	var carts []cart.Cart

	repo.Database.
		Table("carts").
		Where("deleted_at is null").
		Where("phone = ?", phone).
		Order("id asc").
		Scan(&carts)
	return carts
}

package product

import (
	"order-api/pkg/db"

	"gorm.io/gorm/clause"
)

type ProductRepository struct {
	Database *db.Db
}

func NewProductRepository(database *db.Db) *ProductRepository {
	return &ProductRepository{
		Database: database,
	}
}

func (repo *ProductRepository) Create(prod *Product) (*Product, error) {
	result := repo.Database.Create(prod)
	if result.Error != nil {
		return nil, result.Error
	}
	return prod, nil
}

func (repo *ProductRepository) Update(prod *Product) (*Product, error) {
	result := repo.Database.Clauses(clause.Returning{}).Updates(prod)
	if result.Error != nil {
		return nil, result.Error
	}
	return prod, nil
}

func (repo *ProductRepository) Delete(id uint) (error) {
	result := repo.Database.Delete(&Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *ProductRepository) FindById(id uint) (*Product, error) {
	var prod *Product
	result := repo.Database.First(prod, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return prod, nil
}

func (repo *ProductRepository) GetAll(limit, offset int) []Product {
	var products []Product

	repo.Database.
		Table("products").
		Where("deleted_at is null").
		Order("id asc").
		Limit(limit).
		Offset(offset). 
		Scan(&products)
	return products
}
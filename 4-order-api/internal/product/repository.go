package product

import (
	"gorm.io/gorm/clause"
	"order-api/pkg/db"
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
	result := repo.Database.DB.Clauses(clause.Returning{}).Updates(prod)
	if result.Error != nil {
		return nil, result.Error
	}
	return prod, nil
}

func (repo *ProductRepository) Delete(id uint64) error {
	result := repo.Database.Delete(&Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *ProductRepository) FindById(id uint64) (*Product, error) {
	prod := &Product{}
	result := repo.Database.First(prod, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return prod, nil
}

func (repo *ProductRepository) FindByIdUnscoped(id uint64) (*Product, error) {
	prod := &Product{}
	result := repo.Database.Unscoped().Find(prod, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return prod, nil
}

package product

import (
	"order-api/pkg/db"
	"strconv"

	"github.com/lib/pq"
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
	result := repo.Database.DB.Clauses(clause.Returning{}).Updates(prod)
	if result.Error != nil {
		return nil, result.Error
	}
	return prod, nil
}

func (repo *ProductRepository) Delete(id uint64) (error) {
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

// должно быть в сервисе
func (repo *ProductRepository) FindByIDs(cart pq.StringArray) ([]Product, error) {
	products := make([]Product, len(cart))
	for i, productId := range cart {
		prodId, err := strconv.Atoi(productId)
		if err != nil {
			return nil, err
		}
		
		product, err := repo.FindById(uint64(prodId))
		if err != nil {
			return nil, err
		}
		products[i] = *product
	}
	return products, nil
}
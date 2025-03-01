package product

import (
	"fmt"
	"order-api/pkg/db"
	"order-api/pkg/er"

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

// FindByIds находит продукты без метки удалено
func (repo *ProductRepository) FindByIds(ids pq.Int64Array) ([]Product, error) {
	// Поиск всех продуктов
	products := make([]Product, len(ids))

	// Проверяем существуют ли все выбранные продукты
	for i, id := range ids {
		prod, err := repo.FindById(uint64(id))
		if err != nil {
			return nil, er.Wrap(fmt.Sprintf("product №%d", id), err)
		}
		products[i] = *prod
	}
	return products, nil
}

func (repo *ProductRepository) FindByIdUnscoped(id uint64) (*Product, error) {
	prod := &Product{}
	result := repo.Database.Unscoped().Find(prod, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return prod, nil
}


//-------------------------MARKS----------------------------------

func (repo *ProductRepository) AddMark(products []Product, idToAdd uint64) (error) {
	// Обновляем продукты, добавляя id заказа в продукт
	for _, prod := range products {
		prod.CartAdd(int64(idToAdd))
		_, err := repo.Update(&prod)
		if err != nil {
			return err
		}
	}
	return nil
}


func (repo *ProductRepository) DeleteMark(products []Product, idToDelete uint64) (error) { 
	// Логика удаления меток заказа из продуктов
	for _, prod := range products {
		prod.CartRemove(int64(idToDelete))
		_, err := repo.Update(&prod)
		if err != nil {
			return err
		}
	}
	return nil
}

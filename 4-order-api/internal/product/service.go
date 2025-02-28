package product

import (
	"errors"
	"order-api/pkg/er"

	"github.com/lib/pq"
)

type ProductService struct {
	ProductRepository *ProductRepository
}

type ProductServiceDeps struct {
	ProductRepository *ProductRepository
}

func NewProductService(deps ProductServiceDeps) *ProductService {
	return &ProductService{
		ProductRepository: deps.ProductRepository,
	}
}

func (service *ProductService) Create(prod *Product) (*Product, error) {
	return service.ProductRepository.Create(prod)
}

func (service *ProductService) GetByID(id uint64) (*Product, error) {
	return service.ProductRepository.FindById(id)
}

func (service *ProductService) Update(phone string, prod *Product) (*Product, error) {
	if prod.Owner != phone {
		return nil, errors.New(er.ErrWrongUserCredentials)
	}
	return service.ProductRepository.Update(prod)
}

func (service *ProductService) Delete(owner, user string, id uint64) (error) {
	if owner != user {
		return errors.New(er.ErrWrongUserCredentials)
	}
	return service.ProductRepository.Delete(id)
}

func (service *ProductService) GetByIDs(cart pq.StringArray) ([]Product, error) {
	return service.ProductRepository.FindByIDs(cart)
}
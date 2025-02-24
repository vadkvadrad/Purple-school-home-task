package cart

import (
	"errors"
	"fmt"
	"order-api/internal/product"
	"order-api/pkg/er"
	"strconv"
)

type CartService struct {
	CartRepository    *CartRepository
	ProductRepository *product.ProductRepository
}

type CartServiceDeps struct {
	CartRepository    *CartRepository
	ProductRepository *product.ProductRepository
}

func NewCartService(deps CartServiceDeps) *CartService {
	return &CartService{
		CartRepository: deps.CartRepository,
		ProductRepository: deps.ProductRepository,
	}
}

func (service *CartService) Create(cart *Cart) (*Cart, error) {
	for _, id := range cart.Products {
		val, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		_, err = service.ProductRepository.FindById(uint(val))
		if err != nil {
			return nil, er.Wrap(fmt.Sprintf("product №%s", id), err)
		}
	}
	return service.CartRepository.Create(cart)
}

func (service *CartService) GetByIDAndPhone(id uint64, phone string) (*Cart, error) {
	cart, err := service.CartRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if cart.Phone != phone {
		return nil, errors.New(ErrWrongUserCredentials)
	}
	return cart, nil
}

func (service *CartService) GetAll(limit, offset int) []Cart {
	return service.CartRepository.GetAll(limit, offset)
}
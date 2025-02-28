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
		_, err = service.ProductRepository.FindById(uint64(val))
		if err != nil {
			return nil, er.Wrap(fmt.Sprintf("product №%s", id), err)
		}
	}
	return service.CartRepository.Create(cart)
}

func (service *CartService) GetByID(id uint64, phone string) (*Cart, error) {
	cart, err := service.CartRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if cart.Phone != phone {
		return nil, errors.New(er.ErrWrongUserCredentials)
	}
	return cart, nil
}

func (service *CartService) GetByPhone(phone string) []Cart {
	return service.CartRepository.FindByPhone(phone)
}

func (service *CartService) Update(id uint64, newCart *Cart) (*Cart, error) {
	cart, err := service.CartRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if cart.Phone != newCart.Phone {
		return nil, errors.New(er.ErrWrongUserCredentials)
	}
	cart, err = service.CartRepository.Update(newCart)
	if err != nil {
		return nil, err
	}
	return cart, nil
}

func (service *CartService) Delete(id uint64, phone string) (*Cart, error) {
	cart, err := service.CartRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if cart.Phone != phone {
		return nil, errors.New(er.ErrWrongUserCredentials)
	}
	err = service.CartRepository.Delete(id)
	if err != nil {
		return nil, err
	}
	return cart, nil
}
package cart

import (
	"errors"
	"fmt"
	"order-api/internal/product"
	"order-api/pkg/er"
	"strconv"

	"gorm.io/gorm"
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
		CartRepository:    deps.CartRepository,
		ProductRepository: deps.ProductRepository,
	}
}

func (service *CartService) Create(cart *Cart) (*Cart, error) {
	ids := make([]uint, len(cart.Products))
	products := make([]product.Product, len(cart.Products))

	// Проверяем существуют ли все выбранные продукты
	for i, id := range cart.Products {
		val, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		prod, err := service.ProductRepository.FindById(uint64(val))
		if err != nil {
			return nil, er.Wrap(fmt.Sprintf("product №%s", id), err)
		}
		ids[i] = uint(val)
		products[i] = *prod
	}

	// Создаем запись в базе данных
	createdCart, err := service.CartRepository.Create(cart)
	if err != nil {
		return nil, err
	}

	// Обновляем продукты, добавляя id заказа в продукт
	for i, id := range ids {
		prod := products[i]
		carts := prod.Carts
		carts = append(carts, int64(createdCart.ID))
		service.ProductRepository.Update(&product.Product{
			Model:       gorm.Model{ID: id},
			Name:        prod.Name,
			Description: prod.Description,
			Images:      prod.Images,
			Price:       prod.Price,
			Currency:    prod.Currency,
			Owner:       prod.Owner,
			Carts:       carts,
		})
	}
	return createdCart, nil
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

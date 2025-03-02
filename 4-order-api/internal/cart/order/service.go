package order

import (
	"errors"
	"order-api/internal/cart"
	"order-api/internal/product"
	"order-api/pkg/er"
)

type CartService struct {
	CartRepository    cart.ICartRepository
	ProductRepository product.IProductRepository
}

type CartServiceDeps struct {
	CartRepository    cart.ICartRepository
	ProductRepository product.IProductRepository
}

func NewCartService(deps CartServiceDeps) *CartService {
	return &CartService{
		CartRepository:    deps.CartRepository,
		ProductRepository: deps.ProductRepository,
	}
}

func (service *CartService) Create(cart *cart.Cart) (*cart.Cart, error) {
	// Проверяем существуют ли все выбранные продукты
	products, err := service.ProductRepository.FindByIds(cart.Products)
	if err != nil {
		return nil, err
	}

	// Создаем запись в базе данных
	createdCart, err := service.CartRepository.Create(cart)
	if err != nil {
		return nil, err
	}

	// Создаем метку у продуктов
	err = service.ProductRepository.AddMark(products, uint64(createdCart.ID))
	if err != nil {
		return nil, err
	}

	return createdCart, nil
}

func (service *CartService) GetByID(id uint64, phone string) (*cart.Cart, error) {
	cart, err := service.CartRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if cart.Phone != phone {
		return nil, errors.New(er.ErrWrongUserCredentials)
	}
	return cart, nil
}

func (service *CartService) GetByPhone(phone string) []cart.Cart {
	return service.CartRepository.FindByPhone(phone)
}

func (service *CartService) Update(id uint64, newCart *cart.Cart) (*cart.Cart, error) {
	// Поиск заказа
	cart, err := service.CartRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if cart.Phone != newCart.Phone {
		return nil, errors.New(er.ErrWrongUserCredentials)
	}

	// Получение старых продуктов
	oldProducts, err := service.ProductRepository.FindByIds(cart.Products)
	if err != nil {
		return nil, err
	}

	// Получение новых продуктов
	newProducts, err := service.ProductRepository.FindByIds(newCart.Products)
	if err != nil {
		return nil, err
	}

	// Обновление продуктов у заказа
	cart, err = service.CartRepository.Update(newCart)
	if err != nil {
		return nil, err
	}

	// Удаление старых меток у продуктов
	err = service.ProductRepository.DeleteMark(oldProducts, id)
	if err != nil {
		return nil, err
	}

	// Добавление новых меток продуктам
	err = service.ProductRepository.AddMark(newProducts, id)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (service *CartService) Delete(id uint64, phone string) (*cart.Cart, error) {
	// Логика поиска и удаления заказа
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

	// Получение всех продуктов
	products, err := service.ProductRepository.FindByIds(cart.Products)
	if err != nil {
		return nil, err
	}

	// Удаление метки в продукте
	err = service.ProductRepository.DeleteMark(products, id)
	if err != nil {
		return nil, err
	}
	return cart, nil
}

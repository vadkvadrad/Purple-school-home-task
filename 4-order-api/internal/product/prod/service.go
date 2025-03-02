package prod

import (
	"errors"
	"order-api/internal/cart"
	"order-api/internal/product"
	"order-api/internal/user"
	"order-api/pkg/er"
	"order-api/pkg/sender"

	"github.com/lib/pq"
)

type ProductService struct {
	ProductRepository product.IProductRepository
	CartRepository    cart.ICartRepository
	UserRepository    *user.UserRepository
	Sender            *sender.Sender
}

type ProductServiceDeps struct {
	ProductRepository product.IProductRepository
	UserRepository    *user.UserRepository
	CartRepository    cart.ICartRepository
	Sender            *sender.Sender
}

func NewProductService(deps ProductServiceDeps) *ProductService {
	return &ProductService{
		ProductRepository: deps.ProductRepository,
		UserRepository:    deps.UserRepository,
		Sender:            deps.Sender,
	}
}

func (service *ProductService) Create(prod *product.Product) (*product.Product, error) {
	return service.ProductRepository.Create(prod)
}

func (service *ProductService) GetByID(id uint64) (*product.Product, error) {
	return service.ProductRepository.FindById(id)
}

func (service *ProductService) Update(phone string, prod *product.Product) (*product.Product, error) {
	if prod.Owner != phone {
		return nil, errors.New(er.ErrWrongUserCredentials)
	}

	// Обновление продукта
	_, err := service.ProductRepository.Update(prod)
	if err != nil {
		return nil, err
	}

	// // Получение заказов с этим товаром
	// ids := prod.Carts
	// carts := []cart.Cart{}
	// for _, id := range ids {
	// 	service.
	// }

	// // Получение email каждого из заказов
	// emails := []string{}

	// // Отправка письма о изменении заказа
	// text := "Обратите внимание, данные товара " + prod.Name + "были изменены. Зайдите в личный кабинет, чтобы ознакомиться с изменениями."
	// if user.Email != "" {
	// 	err = service.Sender.Email(user.Email, "Товар был изменен", text)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	return prod, nil
}

func (service *ProductService) Delete(owner, user string, id uint64) error {
	if owner != user {
		return errors.New(er.ErrWrongUserCredentials)
	}
	return service.ProductRepository.Delete(id)
}

// GetByIDs находит все продукты, даже с меткой удалено
func (service *ProductService) GetByIDs(cart pq.Int64Array) ([]product.Product, error) {
	products := make([]product.Product, len(cart))
	for i, id := range cart {
		product, err := service.ProductRepository.FindByIdUnscoped(uint64(id))
		if err != nil {
			return nil, err
		}
		products[i] = *product
	}
	return products, nil
}

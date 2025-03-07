package prod

import (
	"errors"
	"order-api/internal/cart"
	"order-api/internal/product"
	"order-api/internal/user"
	"order-api/pkg/er"
	"order-api/pkg/event"
	"order-api/pkg/sender"

	"github.com/lib/pq"
)

type ProductService struct {
	ProductRepository product.IProductRepository
	CartRepository    cart.ICartRepository
	UserRepository    *user.UserRepository
	Sender            *sender.Sender
	EventBus          *event.EventBus
}

type ProductServiceDeps struct {
	ProductRepository product.IProductRepository
	UserRepository    *user.UserRepository
	CartRepository    cart.ICartRepository
	Sender            *sender.Sender
	EventBus          *event.EventBus
}

func NewProductService(deps ProductServiceDeps) *ProductService {
	return &ProductService{
		ProductRepository: deps.ProductRepository,
		CartRepository:    deps.CartRepository,
		UserRepository:    deps.UserRepository,
		Sender:            deps.Sender,
		EventBus:          deps.EventBus,
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

	oldProd, err := service.ProductRepository.FindById(uint64(prod.ID))
	if err != nil {
		return nil, err
	}

	// Обновление продукта
	updatedProd, err := service.ProductRepository.Update(prod)
	if err != nil {
		return nil, err
	}


	// Получение заказов с этим товаром
	ids := updatedProd.Carts
	for _, id := range ids {
		// Получение заказов
		cart, err := service.CartRepository.FindByID(uint64(id))
		if err != nil {
			return nil, err
		}

		// Получение данных пользователя
		user, err := service.UserRepository.FindByKey(user.PhoneKey, cart.Phone)
		if err != nil {
			return nil, err
		}

		// Отправка письма о изменении заказа
		text := "Обратите внимание, данные товара '" + oldProd.Name + "' были изменены. Зайдите в личный кабинет, чтобы ознакомиться с изменениями."
		go service.EventBus.Publish(event.Event{
			Type: event.EventSendEmail,
			Data: sender.Addressee{
				To: user.Email,
				Subject: "Товар был изменен",
				Text: text,
			},
		})
	}

	return prod, nil
}

func (service *ProductService) Delete(owner, buyer string, id uint64) error {
	if owner != buyer {
		return errors.New(er.ErrWrongUserCredentials)
	}

	prod, err := service.ProductRepository.FindById(id)
	if err != nil {
		return err
	}

	// Обновление продукта
	err = service.ProductRepository.Delete(id)
	if err != nil {
		return err
	}

	
	// Получение заказов с этим товаром
	ids := prod.Carts
	for _, id := range ids {
		// Получение заказов
		cart, err := service.CartRepository.FindByID(uint64(id))
		if err != nil {
			return err
		}

		// Получение данных пользователя
		user, err := service.UserRepository.FindByKey(user.PhoneKey, cart.Phone)
		if err != nil {
			return err
		}

		if user.Email == "" {
			continue
		}

		text := "Обратите внимание, товара '" + prod.Name + "' больше нет в наличии. Зайдите в личный кабинет, чтобы ознакомиться с изменениями."
		go service.EventBus.Publish(event.Event{
			Type: event.EventSendEmail,
			Data: sender.Addressee{
				To: user.Email,
				Subject: "Товара нет в наличии",
				Text: text,
			},
		})
	}
	return nil
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

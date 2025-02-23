package cart

import "order-api/internal/product"

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
	for id := range cart.Products {
		_, err := service.ProductRepository.FindById(uint(id))
		if err != nil {
			return nil, err
		}
	}
	return service.CartRepository.Create(cart)
}

func (service *CartService) GetByPhone(phone string) (*Cart, error) {
	return service.CartRepository.FindByPhone(phone)
}

func (service *CartService) GetAll(limit, offset int) []Cart {
	return service.CartRepository.GetAll(limit, offset)
}
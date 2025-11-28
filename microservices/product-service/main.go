package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"product-service/proto/product"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/lib/pq"
	"time"
)

type server struct {
	product.UnimplementedProductServiceServer
	db *gorm.DB
}

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Images      pq.StringArray `json:"images" gorm:"type:text"`
	Price       int32          `json:"price"`
	Currency    string         `json:"currency"`
	Owner       string         `json:"owner"`
	Carts       pq.Int64Array  `json:"carts" gorm:"type:text"`
}

func (s *server) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	newProduct := &Product{
		Name:        req.Name,
		Description: req.Description,
		Images:      pq.StringArray(req.Images),
		Price:       req.Price,
		Currency:    req.Currency,
		Owner:       req.Owner,
		Carts:       pq.Int64Array{},
	}

	result := s.db.Create(newProduct)
	if result.Error != nil {
		return nil, result.Error
	}

	return &product.CreateProductResponse{
		Id:          uint64(newProduct.ID),
		Name:        newProduct.Name,
		Description: newProduct.Description,
		Images:      newProduct.Images,
		Price:       newProduct.Price,
		Currency:    newProduct.Currency,
		Owner:       newProduct.Owner,
		Carts:       newProduct.Carts,
	}, nil
}

func (s *server) GetProduct(ctx context.Context, req *product.GetProductRequest) (*product.GetProductResponse, error) {
	var product Product
	result := s.db.First(&product, "id = ?", req.Id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &product.GetProductResponse{
		Id:          uint64(product.ID),
		Name:        product.Name,
		Description: product.Description,
		Images:      product.Images,
		Price:       product.Price,
		Currency:    product.Currency,
		Owner:       product.Owner,
		Carts:       product.Carts,
	}, nil
}

func (s *server) UpdateProduct(ctx context.Context, req *product.UpdateProductRequest) (*product.UpdateProductResponse, error) {
	var product Product
	result := s.db.First(&product, "id = ?", req.Id)
	if result.Error != nil {
		return nil, result.Error
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Images = pq.StringArray(req.Images)
	product.Price = req.Price
	product.Currency = req.Currency
	product.Owner = req.Owner

	result = s.db.Save(&product)
	if result.Error != nil {
		return nil, result.Error
	}

	return &product.UpdateProductResponse{
		Id:          uint64(product.ID),
		Name:        product.Name,
		Description: product.Description,
		Images:      product.Images,
		Price:       product.Price,
		Currency:    product.Currency,
		Owner:       product.Owner,
		Carts:       product.Carts,
	}, nil
}

func (s *server) DeleteProduct(ctx context.Context, req *product.DeleteProductRequest) (*product.DeleteProductResponse, error) {
	result := s.db.Where("id = ?", req.Id).Delete(&Product{})
	if result.Error != nil {
		return nil, result.Error
	}

	return &product.DeleteProductResponse{
		Success: true,
	}, nil
}

func (s *server) GetProductsByIds(ctx context.Context, req *product.GetProductsByIdsRequest) (*product.GetProductsByIdsResponse, error) {
	var products []Product
	result := s.db.Where("id IN ?", req.Ids).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	var productDataList []*product.ProductData
	for _, prod := range products {
		productDataList = append(productDataList, &product.ProductData{
			Id:          uint64(prod.ID),
			Name:        prod.Name,
			Description: prod.Description,
			Images:      prod.Images,
			Price:       prod.Price,
			Currency:    prod.Currency,
			Owner:       prod.Owner,
			Carts:       prod.Carts,
		})
	}

	return &product.GetProductsByIdsResponse{
		Products: productDataList,
	}, nil
}

func (s *server) AddToCart(ctx context.Context, req *product.AddToCartRequest) (*product.AddToCartResponse, error) {
	var product Product
	result := s.db.First(&product, "id = ?", req.ProductId)
	if result.Error != nil {
		return nil, result.Error
	}

	// Check if cart ID is already in the list
	for _, cartId := range product.Carts {
		if uint64(cartId) == req.CartId {
			return &product.AddToCartResponse{Success: true}, nil
		}
	}

	product.Carts = append(product.Carts, int64(req.CartId))

	result = s.db.Save(&product)
	if result.Error != nil {
		return nil, result.Error
	}

	return &product.AddToCartResponse{
		Success: true,
	}, nil
}

func (s *server) RemoveFromCart(ctx context.Context, req *product.RemoveFromCartRequest) (*product.RemoveFromCartResponse, error) {
	var product Product
	result := s.db.First(&product, "id = ?", req.ProductId)
	if result.Error != nil {
		return nil, result.Error
	}

	var newCarts []int64
	for _, cartId := range product.Carts {
		if uint64(cartId) != req.CartId {
			newCarts = append(newCarts, cartId)
		}
	}
	product.Carts = pq.Int64Array(newCarts)

	result = s.db.Save(&product)
	if result.Error != nil {
		return nil, result.Error
	}

	return &product.RemoveFromCartResponse{
		Success: true,
	}, nil
}

func main() {
	// Connect to database
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=product_db port=5432 sslmode=disable"
	}
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate the schema
	db.AutoMigrate(&Product{})

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	product.RegisterProductServiceServer(s, &server{db: db})
	
	fmt.Println("Product Service is running on port 50052...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"cart-service/proto/cart"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"time"
)

type server struct {
	cart.UnimplementedCartServiceServer
	db *gorm.DB
}

type Cart struct {
	ID       uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Phone    string         `json:"phone"`
	Products pq.Int64Array  `json:"products" gorm:"type:text"`
	Date     datatypes.Date `json:"date"`
}

func (s *server) CreateCart(ctx context.Context, req *cart.CreateCartRequest) (*cart.CreateCartResponse, error) {
	newCart := &Cart{
		Phone:    req.Phone,
		Products: pq.Int64Array(req.Products),
		Date:     datatypes.Date(time.Now()),
	}

	result := s.db.Create(newCart)
	if result.Error != nil {
		return nil, result.Error
	}

	return &cart.CreateCartResponse{
		Id:       uint64(newCart.ID),
		Phone:    newCart.Phone,
		Products: newCart.Products,
		Date:     newCart.Date.String(),
	}, nil
}

func (s *server) GetCart(ctx context.Context, req *cart.GetCartRequest) (*cart.GetCartResponse, error) {
	var cart Cart
	result := s.db.First(&cart, "id = ? AND phone = ?", req.Id, req.Phone)
	if result.Error != nil {
		return nil, result.Error
	}

	return &cart.GetCartResponse{
		Id:       uint64(cart.ID),
		Phone:    cart.Phone,
		Products: cart.Products,
		Date:     cart.Date.String(),
	}, nil
}

func (s *server) UpdateCart(ctx context.Context, req *cart.UpdateCartRequest) (*cart.UpdateCartResponse, error) {
	var cart Cart
	result := s.db.First(&cart, "id = ? AND phone = ?", req.Id, req.Phone)
	if result.Error != nil {
		return nil, result.Error
	}

	cart.Phone = req.Phone
	cart.Products = pq.Int64Array(req.Products)
	cart.Date = datatypes.Date(time.Now())

	result = s.db.Save(&cart)
	if result.Error != nil {
		return nil, result.Error
	}

	return &cart.UpdateCartResponse{
		Id:       uint64(cart.ID),
		Phone:    cart.Phone,
		Products: cart.Products,
		Date:     cart.Date.String(),
	}, nil
}

func (s *server) DeleteCart(ctx context.Context, req *cart.DeleteCartRequest) (*cart.DeleteCartResponse, error) {
	result := s.db.Where("id = ? AND phone = ?", req.Id, req.Phone).Delete(&Cart{})
	if result.Error != nil {
		return nil, result.Error
	}

	return &cart.DeleteCartResponse{
		Success: true,
	}, nil
}

func (s *server) GetCartByPhone(ctx context.Context, req *cart.GetCartByPhoneRequest) (*cart.GetCartByPhoneResponse, error) {
	var carts []Cart
	result := s.db.Where("phone = ?", req.Phone).Find(&carts)
	if result.Error != nil {
		return nil, result.Error
	}

	var cartDataList []*cart.CartData
	for _, cart := range carts {
		cartDataList = append(cartDataList, &cart.CartData{
			Id:       uint64(cart.ID),
			Phone:    cart.Phone,
			Products: cart.Products,
			Date:     cart.Date.String(),
		})
	}

	return &cart.GetCartByPhoneResponse{
		Carts: cartDataList,
	}, nil
}

func main() {
	// Connect to database
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=cart_db port=5432 sslmode=disable"
	}
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate the schema
	db.AutoMigrate(&Cart{})

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	cart.RegisterCartServiceServer(s, &server{db: db})
	
	fmt.Println("Cart Service is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
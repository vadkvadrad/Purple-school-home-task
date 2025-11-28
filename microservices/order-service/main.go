package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"
	"order-service/proto/order"
	productpb "order-service/proto/product"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/lib/pq"
)

type server struct {
	order.UnimplementedOrderServiceServer
	db *gorm.DB
}

type Order struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Phone      string         `json:"phone"`
	ProductIds pq.Int64Array  `json:"product_ids" gorm:"type:text"`
	Status     string         `json:"status"`
	TotalPrice int32          `json:"total_price"`
}

// Product represents the product model for internal use
type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Images      pq.StringArray `json:"images" gorm:"type:text"`
	Price       int32          `json:"price"`
	Currency    string         `json:"currency"`
	Owner       string         `json:"owner"`
	Carts       pq.Int64Array  `json:"carts" gorm:"type:text"`
}

func (s *server) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	newOrder := &Order{
		Phone:      req.Phone,
		ProductIds: pq.Int64Array(req.ProductIds),
		Status:     req.Status,
		TotalPrice: req.TotalPrice,
	}

	result := s.db.Create(newOrder)
	if result.Error != nil {
		return nil, result.Error
	}

	return &order.CreateOrderResponse{
		Id:         uint64(newOrder.ID),
		Phone:      newOrder.Phone,
		ProductIds: newOrder.ProductIds,
		Status:     newOrder.Status,
		TotalPrice: newOrder.TotalPrice,
		CreatedAt:  newOrder.CreatedAt.String(),
	}, nil
}

func (s *server) GetOrder(ctx context.Context, req *order.GetOrderRequest) (*order.GetOrderResponse, error) {
	var order Order
	result := s.db.First(&order, "id = ?", req.Id)
	if result.Error != nil {
		return nil, result.Error
	}

	// Get product details from product service
	productConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer productConn.Close()

	productClient := productpb.NewProductServiceClient(productConn)
	productResp, err := productClient.GetProductsByIds(context.Background(), &productpb.GetProductsByIdsRequest{
		Ids: order.ProductIds,
	})
	if err != nil {
		return nil, err
	}

	return &order.GetOrderResponse{
		Id:         uint64(order.ID),
		Phone:      order.Phone,
		Products:   productResp.Products,
		Status:     order.Status,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt.String(),
	}, nil
}

func (s *server) UpdateOrder(ctx context.Context, req *order.UpdateOrderRequest) (*order.UpdateOrderResponse, error) {
	var order Order
	result := s.db.First(&order, "id = ?", req.Id)
	if result.Error != nil {
		return nil, result.Error
	}

	order.Status = req.Status

	result = s.db.Save(&order)
	if result.Error != nil {
		return nil, result.Error
	}

	// Get product details from product service
	productConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer productConn.Close()

	productClient := productpb.NewProductServiceClient(productConn)
	productResp, err := productClient.GetProductsByIds(context.Background(), &productpb.GetProductsByIdsRequest{
		Ids: order.ProductIds,
	})
	if err != nil {
		return nil, err
	}

	return &order.UpdateOrderResponse{
		Id:         uint64(order.ID),
		Phone:      order.Phone,
		Products:   productResp.Products,
		Status:     order.Status,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt.String(),
	}, nil
}

func (s *server) DeleteOrder(ctx context.Context, req *order.DeleteOrderRequest) (*order.DeleteOrderResponse, error) {
	result := s.db.Where("id = ?", req.Id).Delete(&Order{})
	if result.Error != nil {
		return nil, result.Error
	}

	return &order.DeleteOrderResponse{
		Success: true,
	}, nil
}

func (s *server) GetOrderByPhone(ctx context.Context, req *order.GetOrderByPhoneRequest) (*order.GetOrderByPhoneResponse, error) {
	var orders []Order
	result := s.db.Where("phone = ?", req.Phone).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}

	var orderDataList []*order.OrderData
	for _, ord := range orders {
		// Get product details from product service
		productConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		defer productConn.Close()

		productClient := productpb.NewProductServiceClient(productConn)
		productResp, err := productClient.GetProductsByIds(context.Background(), &productpb.GetProductsByIdsRequest{
			Ids: ord.ProductIds,
		})
		if err != nil {
			return nil, err
		}

		orderDataList = append(orderDataList, &order.OrderData{
			Id:         uint64(ord.ID),
			Phone:      ord.Phone,
			Products:   productResp.Products,
			Status:     ord.Status,
			TotalPrice: ord.TotalPrice,
			CreatedAt:  ord.CreatedAt.String(),
		})
	}

	return &order.GetOrderByPhoneResponse{
		Orders: orderDataList,
	}, nil
}

func (s *server) ProcessCartToOrder(ctx context.Context, req *order.ProcessCartToOrderRequest) (*order.ProcessCartToOrderResponse, error) {
	// In a real implementation, you would:
	// 1. Call cart service to get cart details
	// 2. Call product service to get product details and calculate total price
	// 3. Create an order with the cart data
	// 4. Update cart status or delete the cart
	// 5. Return the new order ID

	// For this example, we'll create a simple order
	newOrder := &Order{
		Phone:      req.Phone,
		ProductIds: pq.Int64Array{}, // This would be filled from cart data
		Status:     "created",
		TotalPrice: 0, // This would be calculated from product prices
	}

	result := s.db.Create(newOrder)
	if result.Error != nil {
		return nil, result.Error
	}

	return &order.ProcessCartToOrderResponse{
		OrderId: uint64(newOrder.ID),
		Status:  "created",
		Success: true,
	}, nil
}

func main() {
	// Connect to database
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=order_db port=5432 sslmode=disable"
	}
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate the schema
	db.AutoMigrate(&Order{})

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	order.RegisterOrderServiceServer(s, &server{db: db})
	
	fmt.Println("Order Service is running on port 50053...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
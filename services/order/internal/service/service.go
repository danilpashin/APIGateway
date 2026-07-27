package service

import (
	"apigateway/services/order/internal/domain"
	"apigateway/services/order/internal/repository/postgres"
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type OrderService struct {
	repo postgres.OrderRepositoryInterface
}

func NewOrderService(repo postgres.OrderRepositoryInterface) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) AddProductToCart(ctx context.Context, req *domain.AddProductToCartRequest) (*domain.Order, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	key := fmt.Sprintf("cart_%d", req.UserID)

	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		fmt.Printf("key `%s` does not exist", key)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get products from %s, err: %v", key, err)
	}

	fmt.Println("val: ", val)

	err = rdb.Set(ctx, key, "value", 123).Err()
	if err != nil {
		return nil, errors.New("failed to add product to cart")
	}

	return nil, nil
}

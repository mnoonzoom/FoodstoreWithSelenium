package service

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"order/internal/dao"
	"order/internal/model"
	"time"
)

type OrderService struct {
	repo *dao.OrderDao
}

func NewOrderService(repo *dao.OrderDao) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID string, itemIDs []string, totalPrice float64) (string, error) {
	order := model.Order{
		UserID:     userID,
		ItemIDs:    itemIDs,
		TotalPrice: totalPrice,
		Status:     "Pending",
		CreatedAt:  time.Now(),
	}
	return s.repo.Create(ctx, order)
}

func (s *OrderService) GetOrder(ctx context.Context, id string) (*model.Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, id string, status string) error {
	return s.repo.UpdateStatus(ctx, id, status)
}
func (s *OrderService) UpdateOrder(ctx context.Context, order model.Order) error {

	_, err := s.repo.Collection.ReplaceOne(
		ctx,
		bson.M{"_id": order.ID},
		order,
	)
	return err
}

func (s *OrderService) DeleteOrder(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
func (s *OrderService) ListOrdersByUser(ctx context.Context, userId string) ([]model.Order, error) {
	return s.repo.FindOrdersByUserId(ctx, userId)
}

func (s *OrderService) ListOrders(ctx context.Context, limit int64, skip int64) ([]model.Order, error) {
	return s.repo.List(ctx, limit, skip)
}

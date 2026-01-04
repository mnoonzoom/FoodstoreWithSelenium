package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"order/internal/model"
	"order/internal/service"
)

type MockOrderDao struct {
	mock.Mock
}

func (m *MockOrderDao) Create(ctx context.Context, order model.Order) (string, error) {
	args := m.Called(ctx, order)
	return args.String(0), args.Error(1)
}

func (m *MockOrderDao) GetByID(ctx context.Context, id string) (*model.Order, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderDao) UpdateStatus(ctx context.Context, id string, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockOrderDao) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrderDao) FindOrdersByUserId(ctx context.Context, userId string) ([]model.Order, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).([]model.Order), args.Error(1)
}

func (m *MockOrderDao) List(ctx context.Context, limit int64, skip int64) ([]model.Order, error) {
	args := m.Called(ctx, limit, skip)
	return args.Get(0).([]model.Order), args.Error(1)
}

func TestOrderService_CreateOrder(t *testing.T) {
	mockRepo := new(MockOrderDao)
	svc := service.NewOrderService(mockRepo)

	userID := "user123"
	itemIDs := []string{"item1", "item2"}
	totalPrice := 50.0

	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(order model.Order) bool {
		return order.UserID == userID &&
			len(order.ItemIDs) == len(itemIDs) &&
			order.TotalPrice == totalPrice &&
			order.Status == "Pending"
	})).Return("order123", nil)

	id, err := svc.CreateOrder(context.Background(), userID, itemIDs, totalPrice)

	assert.NoError(t, err)
	assert.Equal(t, "order123", id)

	mockRepo.AssertExpectations(t)
}

func TestOrderService_GetOrder(t *testing.T) {
	mockRepo := new(MockOrderDao)
	svc := service.NewOrderService(mockRepo)

	order := &model.Order{ID: "order123", UserID: "user123", Status: "Pending"}

	mockRepo.On("GetByID", mock.Anything, "order123").Return(order, nil)

	res, err := svc.GetOrder(context.Background(), "order123")

	assert.NoError(t, err)
	assert.Equal(t, order, res)

	mockRepo.AssertExpectations(t)
}

func TestOrderService_UpdateOrderStatus(t *testing.T) {
	mockRepo := new(MockOrderDao)
	svc := service.NewOrderService(mockRepo)

	mockRepo.On("UpdateStatus", mock.Anything, "order123", "Completed").Return(nil)

	err := svc.UpdateOrderStatus(context.Background(), "order123", "Completed")

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestOrderService_DeleteOrder(t *testing.T) {
	mockRepo := new(MockOrderDao)
	svc := service.NewOrderService(mockRepo)

	mockRepo.On("Delete", mock.Anything, "order123").Return(nil)

	err := svc.DeleteOrder(context.Background(), "order123")

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestOrderService_ListOrdersByUser(t *testing.T) {
	mockRepo := new(MockOrderDao)
	svc := service.NewOrderService(mockRepo)

	orders := []model.Order{
		{ID: "order1", UserID: "user123"},
		{ID: "order2", UserID: "user123"},
	}

	mockRepo.On("FindOrdersByUserId", mock.Anything, "user123").Return(orders, nil)

	res, err := svc.ListOrdersByUser(context.Background(), "user123")

	assert.NoError(t, err)
	assert.Equal(t, orders, res)

	mockRepo.AssertExpectations(t)
}

func TestOrderService_ListOrders(t *testing.T) {
	mockRepo := new(MockOrderDao)
	svc := service.NewOrderService(mockRepo)

	orders := []model.Order{
		{ID: "order1"},
		{ID: "order2"},
	}

	mockRepo.On("List", mock.Anything, int64(10), int64(0)).Return(orders, nil)

	res, err := svc.ListOrders(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, orders, res)

	mockRepo.AssertExpectations(t)
}

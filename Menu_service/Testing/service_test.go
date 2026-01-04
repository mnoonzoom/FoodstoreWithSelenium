package service_test

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"menu/internal/model"
	"menu/internal/service"
)

type MockMenuRepo struct {
	mock.Mock
}

func (m *MockMenuRepo) CreateMenuItem(ctx context.Context, item model.MenuItem) (string, error) {
	args := m.Called(ctx, item)
	return args.String(0), args.Error(1)
}

func (m *MockMenuRepo) GetAllMenuItems(ctx context.Context, filter interface{}, limit, skip int64, sortBy string, asc bool) ([]model.MenuItem, error) {
	args := m.Called(ctx, filter, limit, skip, sortBy, asc)
	return args.Get(0).([]model.MenuItem), args.Error(1)
}

func (m *MockMenuRepo) GetMenuItemByID(ctx context.Context, id string) (*model.MenuItem, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.MenuItem), args.Error(1)
}

func (m *MockMenuRepo) Update(ctx context.Context, id string, update primitive.M) error {
	args := m.Called(ctx, id, update)
	return args.Error(0)
}

func (m *MockMenuRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMenuRepo) CountMenuItems(ctx context.Context, filter interface{}) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func TestMenuService_AllMethods(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockMenuRepo)
	svc := service.NewMenuService(mockRepo)

	t.Log("Testing CreateMenuItem")
	item := model.MenuItem{Name: "Pizza", Price: 9.99, Category: "Main"}
	mockRepo.On("CreateMenuItem", mock.Anything, item).Return("123", nil)
	id, err := svc.CreateMenuItem(ctx, item)
	t.Logf("CreateMenuItem returned ID: %s, error: %v", id, err)
	assert.NoError(t, err)
	assert.Equal(t, "123", id)
	mockRepo.AssertCalled(t, "CreateMenuItem", mock.Anything, item)

	t.Log("Testing GetAllMenuItems")
	expectedItems := []model.MenuItem{item}
	mockRepo.On("GetAllMenuItems", mock.Anything, bson.M{}, int64(10), int64(0), "", true).Return(expectedItems, nil)
	items, err := svc.GetAllMenuItems(ctx, bson.M{}, 10, 0, "", true)
	t.Logf("GetAllMenuItems returned %d items, error: %v", len(items), err)
	assert.NoError(t, err)
	assert.Equal(t, expectedItems, items)
	mockRepo.AssertCalled(t, "GetAllMenuItems", mock.Anything, bson.M{}, int64(10), int64(0), "", true)

	t.Log("Testing GetMenuItemByID")
	mockRepo.On("GetMenuItemByID", mock.Anything, "123").Return(&item, nil)
	gotItem, err := svc.GetMenuItemByID(ctx, "123")
	t.Logf("GetMenuItemByID returned item: %+v, error: %v", gotItem, err)
	assert.NoError(t, err)
	assert.Equal(t, &item, gotItem)
	mockRepo.AssertCalled(t, "GetMenuItemByID", mock.Anything, "123")

	t.Log("Testing UpdateMenuItem")
	update := primitive.M{"price": 12.99}
	mockRepo.On("Update", mock.Anything, "123", update).Return(nil)
	err = svc.UpdateMenuItem(ctx, "123", update)
	t.Logf("UpdateMenuItem error: %v", err)
	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "Update", mock.Anything, "123", update)

	t.Log("Testing DeleteMenuItem")
	mockRepo.On("Delete", mock.Anything, "123").Return(nil)
	err = svc.DeleteMenuItem(ctx, "123")
	t.Logf("DeleteMenuItem error: %v", err)
	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "Delete", mock.Anything, "123")

	t.Log("Testing CountMenuItems")
	mockRepo.On("CountMenuItems", mock.Anything, bson.M{}).Return(int64(5), nil)
	count, err := svc.CountMenuItems(ctx, bson.M{})
	t.Logf("CountMenuItems returned count: %d, error: %v", count, err)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
	mockRepo.AssertCalled(t, "CountMenuItems", mock.Anything, bson.M{})

	mockRepo.AssertExpectations(t)
}

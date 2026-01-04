package service

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"foodstore/menu/internal/dao"
	"foodstore/menu/internal/model"
)

type MenuService struct {
	repo dao.MenuRepository
}

func NewMenuService(repo dao.MenuRepository) *MenuService {
	return &MenuService{repo: repo}
}

func (s *MenuService) CreateMenuItem(ctx context.Context, item model.MenuItem) (string, error) {
	return s.repo.CreateMenuItem(ctx, item)
}

func (s *MenuService) GetMenuItemByID(ctx context.Context, id string) (*model.MenuItem, error) {
	return s.repo.GetMenuItemByID(ctx, id)
}

func (s *MenuService) UpdateMenuItem(ctx context.Context, id string, update bson.M) error {
	return s.repo.Update(ctx, id, update)
}

func (s *MenuService) DeleteMenuItem(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *MenuService) GetAllMenuItems(ctx context.Context, filter bson.M, limit, skip int64, sortBy string, asc bool) ([]model.MenuItem, error) {
	return s.repo.GetAllMenuItems(ctx, filter, limit, skip, sortBy, asc)
}

func (s *MenuService) GetMultipleMenuItems(ctx context.Context, ids []string) ([]model.MenuItem, error) {
	objectIDs := []interface{}{}
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue
		}
		objectIDs = append(objectIDs, oid)
	}
	filter := bson.M{"_id": bson.M{"$in": objectIDs}}

	return s.repo.GetAllMenuItems(ctx, filter, 0, 0, "", true)
}

func (s *MenuService) CountMenuItems(ctx context.Context, filter bson.M) (int64, error) {
	return s.repo.CountMenuItems(ctx, filter)
}

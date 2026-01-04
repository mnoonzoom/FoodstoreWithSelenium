package handler

import (
	"context"
	"foodstore/menu/internal/model"
	"foodstore/menu/internal/service"
	pb "foodstore/menu/proto"

	"go.mongodb.org/mongo-driver/bson"
)

type MenuHandler struct {
	pb.UnimplementedMenuServiceServer
	menuService *service.MenuService
}

func NewMenuHandler(service *service.MenuService) *MenuHandler {
	return &MenuHandler{menuService: service}
}

func (h *MenuHandler) ListMenuItems(ctx context.Context, req *pb.ListMenuItemsRequest) (*pb.ListMenuItemsResponse, error) {
	filter := bson.M{}

	if req.Category != "" {
		filter["category"] = req.Category
	}

	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": req.Search, "$options": "i"}},
			{"description": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}

	totalCount, err := h.menuService.CountMenuItems(ctx, filter)
	if err != nil {
		return nil, err
	}

	items, err := h.menuService.GetAllMenuItems(ctx, filter, req.Limit, req.Skip, req.SortBy, req.SortAsc)
	if err != nil {
		return nil, err
	}

	var responseItems []*pb.MenuItem
	for _, item := range items {
		responseItems = append(responseItems, &pb.MenuItem{
			Id:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			Available:   item.Available,
			Category:    item.Category,
			ImageUrl:    item.ImageURL,
		})
	}

	return &pb.ListMenuItemsResponse{
		Items:      responseItems,
		TotalCount: totalCount,
	}, nil
}

func (h *MenuHandler) CreateMenuItem(ctx context.Context, req *pb.CreateMenuItemRequest) (*pb.CreateMenuItemResponse, error) {
	item := model.MenuItem{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Available:   req.Available,
		Category:    req.Category,
		ImageURL:    req.ImageUrl,
	}

	id, err := h.menuService.CreateMenuItem(ctx, item)
	if err != nil {
		return nil, err
	}

	return &pb.CreateMenuItemResponse{Id: id}, nil
}

func (h *MenuHandler) GetMenuItemByID(ctx context.Context, req *pb.GetMenuItemByIDRequest) (*pb.GetMenuItemByIDResponse, error) {
	item, err := h.menuService.GetMenuItemByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetMenuItemByIDResponse{
		Item: &pb.MenuItem{
			Id:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			Available:   item.Available,
			Category:    item.Category,
			ImageUrl:    item.ImageURL,
		},
	}, nil
}

func (h *MenuHandler) UpdateMenuItem(ctx context.Context, req *pb.UpdateMenuItemRequest) (*pb.UpdateMenuItemResponse, error) {
	update := bson.M{}

	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Description != "" {
		update["description"] = req.Description
	}
	if req.Price != 0 {
		update["price"] = req.Price
	}
	update["available"] = req.Available
	if req.Category != "" {
		update["category"] = req.Category
	}
	if req.ImageUrl != "" {
		update["image_url"] = req.ImageUrl
	}

	err := h.menuService.UpdateMenuItem(ctx, req.Id, update)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateMenuItemResponse{Message: "Updated successfully"}, nil
}

func (h *MenuHandler) DeleteMenuItem(ctx context.Context, req *pb.DeleteMenuItemRequest) (*pb.DeleteMenuItemResponse, error) {
	err := h.menuService.DeleteMenuItem(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteMenuItemResponse{Message: "Deleted successfully"}, nil
}
func (h *MenuHandler) GetMultipleMenuItems(ctx context.Context, req *pb.GetMultipleMenuItemsRequest) (*pb.GetMultipleMenuItemsResponse, error) {
	items, err := h.menuService.GetMultipleMenuItems(ctx, req.Ids)
	if err != nil {
		return nil, err
	}

	var responseItems []*pb.MenuItem
	for _, item := range items {
		responseItems = append(responseItems, &pb.MenuItem{
			Id:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			Available:   item.Available,
			Category:    item.Category,
			ImageUrl:    item.ImageURL,
		})
	}

	return &pb.GetMultipleMenuItemsResponse{Items: responseItems}, nil
}

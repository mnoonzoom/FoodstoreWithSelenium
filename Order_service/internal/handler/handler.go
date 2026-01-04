package handler

import (
	"context"
	"fmt"
	"order/internal/model"
	nats "order/internal/nats"

	"order/internal/service"
	pb "order/proto"
	menupb "order/proto/menu"
	"time"
)

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer
	svc           *service.OrderService
	menuClient    menupb.MenuServiceClient
	natsPublisher *nats.Publisher
}

func NewOrderHandler(
	svc *service.OrderService,
	menuClient menupb.MenuServiceClient,
	natsPublisher *nats.Publisher,
) *OrderHandler {
	return &OrderHandler{
		svc:           svc,
		menuClient:    menuClient,
		natsPublisher: natsPublisher,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	menuRes, err := h.menuClient.GetMultipleMenuItems(ctx, &menupb.GetMultipleMenuItemsRequest{
		Ids: req.ItemIds,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch menu items: %v", err)
	}

	totalPrice := 0.0
	for _, item := range menuRes.Items {
		totalPrice += item.Price
	}

	id, err := h.svc.CreateOrder(ctx, req.UserId, req.ItemIds, totalPrice)
	if err != nil {
		return nil, err
	}

	_ = h.natsPublisher.PublishOrderCreated(map[string]interface{}{
		"orderId":   id,
		"userId":    req.UserId,
		"items":     req.ItemIds,
		"total":     totalPrice,
		"createdAt": time.Now().Format(time.RFC3339),
	})

	return &pb.CreateOrderResponse{Id: id}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	order, err := h.svc.GetOrder(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetOrderResponse{
		Order: &pb.Order{
			Id:         order.ID,
			UserId:     order.UserID,
			ItemIds:    order.ItemIDs,
			TotalPrice: order.TotalPrice,
			Status:     order.Status,
			CreatedAt:  order.CreatedAt.String(),
		},
	}, nil
}

func (h *OrderHandler) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	order := model.Order{
		ID:         req.Id,
		UserID:     req.UserId,
		ItemIDs:    req.ItemIds,
		TotalPrice: req.TotalPrice,
		Status:     req.Status,
	}

	err := h.svc.UpdateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateOrderResponse{Message: "Order fully updated"}, nil
}
func (h *OrderHandler) PatchOrderStatus(ctx context.Context, req *pb.PatchOrderStatusRequest) (*pb.PatchOrderStatusResponse, error) {
	err := h.svc.UpdateOrderStatus(ctx, req.Id, req.Status)
	if err != nil {
		return nil, err
	}

	return &pb.PatchOrderStatusResponse{Message: "Status updated"}, nil
}

func (h *OrderHandler) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	err := h.svc.DeleteOrder(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteOrderResponse{Message: "Order deleted successfully"}, nil
}

func (h *OrderHandler) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := h.svc.ListOrders(ctx, req.Limit, req.Skip)
	if err != nil {
		return nil, err
	}

	var pbOrders []*pb.Order
	for _, order := range orders {
		pbOrders = append(pbOrders, &pb.Order{
			Id:         order.ID,
			UserId:     order.UserID,
			ItemIds:    order.ItemIDs,
			TotalPrice: order.TotalPrice,
			Status:     order.Status,
			CreatedAt:  order.CreatedAt.String(),
		})
	}

	return &pb.ListOrdersResponse{Orders: pbOrders}, nil
}
func (h *OrderHandler) ListOrdersByUser(ctx context.Context, req *pb.ListOrdersByUserRequest) (*pb.ListOrdersByUserResponse, error) {
	orders, err := h.svc.ListOrdersByUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	var pbOrders []*pb.Order
	for _, order := range orders {
		pbOrders = append(pbOrders, &pb.Order{
			Id:         order.ID,
			UserId:     order.UserID,
			ItemIds:    order.ItemIDs,
			TotalPrice: order.TotalPrice,
			Status:     order.Status,
			CreatedAt:  order.CreatedAt.String(),
		})
	}

	return &pb.ListOrdersByUserResponse{
		Orders: pbOrders,
	}, nil
}

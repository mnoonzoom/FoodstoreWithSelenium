package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"log"
	"net"
	"order/config"
	"order/internal/dao"
	"order/internal/handler"
	"order/internal/nats"
	"order/internal/service"
	pb "order/proto"
	menupb "order/proto/menu"
)

func main() {
	cfg := config.LoadConfig()
	db := config.ConnectToMongo(cfg.MongoURI, cfg.DatabaseName)
	cache := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	repo := dao.NewOrderDao(db, cache)
	svc := service.NewOrderService(repo)

	menuConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to MenuService: %v", err)
	}
	menuClient := menupb.NewMenuServiceClient(menuConn)
	natsPublisher, err := nats.NewPublisher("nats://localhost:4222")
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	orderHandler := handler.NewOrderHandler(svc, menuClient, natsPublisher)

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, orderHandler)

	fmt.Println("OrderService is running on :50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}
}

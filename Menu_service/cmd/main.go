package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"net"

	"foodstore/menu/config"
	"foodstore/menu/internal/dao"
	"foodstore/menu/internal/handler"
	"foodstore/menu/internal/service"
	pb "foodstore/menu/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMenuServiceServer
	menuService *service.MenuService
}

func main() {
	cfg := config.LoadConfig()
	db := config.ConnectToMongo(cfg.MongoURI, cfg.DatabaseName)
	cache := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	menuRepo := dao.NewMenuRepository(db, cache)
	menuService := service.NewMenuService(menuRepo)
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMenuServiceServer(grpcServer, handler.NewMenuHandler(menuService))
	fmt.Println("gRPC server started on port :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}
}

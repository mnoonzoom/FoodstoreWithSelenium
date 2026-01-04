package main

import (
	"fmt"
	"log"
	"net"
	"user/config"
	"user/internal/dao"
	"user/internal/handler"
	"user/internal/service"
	pb "user/proto"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()

	db := config.ConnectToMongo(cfg.MongoURI, cfg.DatabaseName)
	repo := dao.NewUserRepository(db)
	svc := service.NewUserService(repo)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, handler.NewUserHandler(svc))

	fmt.Println("UserService started on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}
}

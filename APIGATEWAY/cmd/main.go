package main

import (
	"apigateway/internal/handler"

	menuPB "apigateway/proto/menu"
	orderPB "apigateway/proto/order"
	userPB "apigateway/proto/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"google.golang.org/grpc"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8082"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	menuConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to MenuService: %v", err)
	}
	orderConn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to OrderService: %v", err)
	}
	userConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to UserService: %v", err)
	}

	menuClient := menuPB.NewMenuServiceClient(menuConn)
	orderClient := orderPB.NewOrderServiceClient(orderConn)
	userClient := userPB.NewUserServiceClient(userConn)

	handler.InitMenuRoutes(r, menuClient)
	handler.InitOrderRoutes(r, orderClient)
	handler.InitUserRoutes(r, userClient)

	log.Println("API Gateway started on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("API Gateway failed: %v", err)
	}
}

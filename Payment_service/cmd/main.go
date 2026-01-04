package main

import (
	"context"
	"github.com/joho/godotenv"
	natslib "github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"log"
	"os"
	"payment/mailer"
	"payment/nats"
	userpb "payment/proto/user"
	"strconv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found, using system env")
	}

	nc, err := natslib.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	userConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to UserService: %v", err)
	}
	defer userConn.Close()

	userClient := userpb.NewUserServiceClient(userConn)
	rawPort := os.Getenv("SMTP_PORT")
	if rawPort == "" {
		log.Fatal("SMTP_PORT is missing in environment")
	}
	port, err := strconv.Atoi(rawPort)
	if err != nil || port <= 0 {
		log.Fatalf("Invalid SMTP_PORT value: %s", rawPort)
	}

	m := mailer.NewMailer(
		os.Getenv("SMTP_HOST"),
		port,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"),
		os.Getenv("SMTP_FROM"),
	)
	worker := &nats.EmailWorker{
		Mailer: m,
		GetEmailFn: func(userID string) (string, error) {
			res, err := userClient.GetUser(context.Background(), &userpb.GetUserRequest{Id: userID})
			if err != nil {
				log.Printf("Failed to fetch user via gRPC: %v", err)
				return "", err
			}
			return res.User.Email, nil
		},
	}

	if _, err := nc.Subscribe("order.created", worker.HandleOrderCreated); err != nil {
		log.Fatal(err)
	}

	log.Println("EmailService is listening on order.created...")
	select {}
}

# QuickBite – Microservices-based Food Store System

## Project Overview and Topic

**QuickBite** is a full-stack food ordering application built with a microservices architecture. It allows users to browse the menu, place orders, and receive receipts via email. The system is divided into multiple services communicating via **gRPC** and **NATS** message queue. Project done by: Adilkhan Dias, Arman Bezhanov, Danial Turzhanov. Group SE-2331

## Technologies Used

- **Go** – for backend microservices
- **gRPC** – inter-service communication
- **MongoDB** – primary database
- **Redis** – optional caching
- **NATS** – event-driven messaging
- **Gin** – for API Gateway
- **gomail** – for sending email
- **JWT** – user authentication
- **JavaScript + HTML/CSS** – frontend (static served with Go)

## How to Run Locally

### Prerequisites

- Go 1.20+
- MongoDB running on `localhost:27017`
- NATS server (`nats-server` running on `localhost:4222`)

### Steps

1. **Start MongoDB and NATS**

   ```bash
   mongod
   nats-server -DV
   ```

2. **Start each service manually:**

   ```bash
   cd User_service/cmd && go run main.go
   cd Menu_service/cmd && go run main.go
   cd Order_service/cmd && go run main.go
   cd Payment_service/cmd && go run main.go
   cd APIGATEWAY && go run main.go
   cd Frontend && go run main.go
   ```

3. Visit `http://localhost:8082` to open the app.

## How to Run Tests

```bash
cd Menu_service/internal/service
go test -v
```

> Repeat for other services with `go test ./...`

## Description of gRPC Endpoints

### MenuService

- `CreateMenuItem(CreateMenuItemRequest) returns (CreateMenuItemResponse)`
- `GetMenuItemByID(GetMenuItemByIDRequest) returns (GetMenuItemByIDResponse)`
- `UpdateMenuItem(UpdateMenuItemRequest) returns (UpdateMenuItemResponse)`
- `DeleteMenuItem(DeleteMenuItemRequest) returns (DeleteMenuItemResponse)`
- `ListMenuItems(ListMenuItemsRequest) returns (ListMenuItemsResponse)`
- `GetMultipleMenuItems(GetMultipleMenuItemsRequest) returns (GetMultipleMenuItemsResponse)`

### OrderService

- `CreateOrder(CreateOrderRequest) returns (CreateOrderResponse)`
- `GetOrder(GetOrderRequest) returns (GetOrderResponse)`
- `UpdateOrder(UpdateOrderRequest) returns (UpdateOrderResponse)`
- `PatchOrderStatus(PatchOrderStatusRequest) returns (PatchOrderStatusResponse)`
- `DeleteOrder(DeleteOrderRequest) returns (DeleteOrderResponse)`
- `ListOrders(ListOrdersRequest) returns (ListOrdersResponse)`

### UserService

- `Register(RegisterRequest) returns (RegisterResponse)`
- `Login(LoginRequest) returns (LoginResponse)`
- `GetUser(GetUserRequest) returns (GetUserResponse)`

## List of Implemented Features

User Registration & Login with JWT

Menu Management (CRUD)

Order Creation & Tracking

NATS-based message event

Sending Email Receipts via Gmail SMTP

API Gateway (Gin)

gRPC communication between services

Unit Tests for services

Static Frontend served via Go

Pagination, Sorting, and Filtering (menu)

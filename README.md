# YouAre — Microservices Backend

YouAre is a backend project built with a microservices architecture using Go. It manages users, products, and orders, and supports asynchronous communication via RabbitMQ. Each service communicates via gRPC and uses MongoDB for persistence.

---

## Project Overview

YouAre is a distributed system designed for managing a simple online store. It supports:

- User registration and authentication
- Product catalog management
- Order placement and stock update
- Event-driven communication via RabbitMQ

---

## Technologies Used

- **Go** — primary language for all services
- **gRPC** — inter-service communication
- **MongoDB** — primary database
- **Redis** — caching layer
- **RabbitMQ** — asynchronous messaging
- **Docker / Docker Compose** — for rabbitMQ and redis
- **Protocol Buffers (protobuf)** — API definition

---

## How to Run Locally

### Option 1: With Docker Compose (Redis and RabbitMQ only)

1. Start RabbitMQ and Redis:
   ```bash
   docker compose up -d
2. Start each service manually in separate terminals:
    cd user-service
    go run cmd/main.go

    cd product-service
    go run cmd/main.go

    cd order-service
    go run cmd/main.go

    cd consumer-service
    go run cmd/main.go

    cd api-gateway
    go run cmd/main.go
3. Make sure to start mongoDB on localhost:27017


## How to Run Tests

### Each service contains unit and integration tests in the internal/.../tests directories.
    
    go test ./...

## gRPC Endpoints

### UserService (user-service)

| Method           | Request Type            | Response Type           | Description                                                         |
| ---------------- | ----------------------- | ----------------------- | ------------------------------------------------------------------- |
| `Register`       | `RegisterRequest`       | `UserResponse`          | Registers a new user with the provided credentials.                 |
| `Login`          | `LoginRequest`          | `LoginResponse`         | Authenticates a user and returns a JWT token or similar login data. |
| `GetUser`        | `GetUserRequest`        | `UserResponse`          | Retrieves user details by ID or email.                              |
| `GetProfile`     | `ProfileRequest`        | `UserResponse`          | Retrieves the authenticated user's profile.                         |
| `GetAllProfiles` | `google.protobuf.Empty` | `UserListResponse`      | Returns a list of all user profiles.                                |
| `DeleteUser`     | `ProfileRequest`        | `google.protobuf.Empty` | Deletes the user profile associated with the given ID or token.     |


### ProductService (product-service)

| Method     | Request Type            | Response Type           | Description                                                                  |
| ---------- | ----------------------- | ----------------------- | ---------------------------------------------------------------------------- |
| `Create`   | `CreateRequest`         | `ProductResponse`       | Creates a new product with details like name, description, price, and stock. |
| `Get`      | `ProductRequest`        | `ProductResponse`       | Retrieves a single product by ID.                                            |
| `GetAll`   | `google.protobuf.Empty` | `ProductListResponse`   | Returns a list of all products.                                              |
| `Update`   | `UpdateRequest`         | `ProductResponse`       | Updates product details.                                                     |
| `Decrease` | `DecreaseRequest`       | `ProductResponse`       | Decreases the stock of a product, typically used after an order.             |
| `Delete`   | `ProductRequest`        | `google.protobuf.Empty` | Deletes a product by ID.                                                     |

## OrderService (order-service)

| Method         | Request Type            | Response Type           | Description                                           |
| -------------- | ----------------------- | ----------------------- | ----------------------------------------------------- |
| `CreateOrder`  | `CreateOrderRequest`    | `OrderResponse`         | Creates a new order with user and product references. |
| `GetOrder`     | `OrderRequest`          | `OrderResponse`         | Retrieves a specific order by ID.                     |
| `GetAllOrders` | `google.protobuf.Empty` | `OrderListResponse`     | Returns a list of all orders.                         |
| `DeleteOrder`  | `OrderRequest`          | `google.protobuf.Empty` | Deletes an order by ID.                               |

## Implemented Features

### User Service
- User registration
- User login and authentication (with JWT)
- Fetch user profile
- Admin: fetch all users
- Delete user profile

### Product Service
- Create a new product
- Get product by ID
- List all products
- Update product info
- Decrease product stock
- Delete a product

### Order Service
- Create an order
- Get order by ID
- List all orders
- Delete order

### Event-Driven Features
- Asynchronous stock update via RabbitMQ after order creation
- Consumer service that listens to order events and updates product stock

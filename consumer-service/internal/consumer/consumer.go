package consumer

import (
	"consumer-service/internal/entity"
	"consumer-service/internal/grpc"
	"consumer-service/internal/logger"
	"consumer-service/internal/repository"
	"encoding/json"
	userpb "user-service/pkg/gen/user"

	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn          *amqp091.Connection
	productClient *grpc.ProductClient
	userClient    *grpc.UserClient
	repo          repository.OrderRepository
	logger        *logger.Logger
}

func NewConsumer(conn *amqp091.Connection, productClient *grpc.ProductClient, userClient *grpc.UserClient, repo repository.OrderRepository) *Consumer {
	return &Consumer{
		conn:          conn,
		productClient: productClient,
		userClient:    userClient,
		repo:          repo,
		logger:        logger.GetLogger(),
	}
}

func (c *Consumer) StartConsuming() error {
	ch, err := c.conn.Channel()
	if err != nil {
		c.logger.Error("Failed to open a channel: %v", err)
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"order.created",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		c.logger.Error("Failed to declare queue: %v", err)
		return err
	}

	msgs, err := ch.Consume(
		"order.created",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		c.logger.Error("Failed to register consumer: %v", err)
		return err
	}

	c.logger.Info(" [*] Successfully connected to RabbitMQ and waiting for messages from order.created...")

	for msg := range msgs {
		c.logger.Info("Received new order message, starting processing")
		c.logger.Debug("Message body: %s", string(msg.Body))

		var order entity.Order
		err := json.Unmarshal(msg.Body, &order)
		if err != nil {
			c.logger.Error("Failed to parse order message: %v. Message will be rejected", err)
			msg.Nack(false, false)
			continue
		}

		c.logger.Info("Processing order ID: %s for user ID: %s and product ID: %s",
			order.ID.Hex(), order.UserID, order.ProductID)

		c.logger.Debug("Validating user ID: %s", order.UserID)
		profileReq := &userpb.ProfileRequest{UserId: order.UserID}
		user, err := c.userClient.GetProfile(profileReq)
		if err != nil {
			c.logger.Error("Failed to validate user ID: %s: %v. Message will be rejected", order.UserID, err)
			msg.Nack(false, false)
			continue
		}
		c.logger.Info("Order placed by valid user: %s (email: %s)", user.Id, user.Email)

		c.logger.Debug("Requesting product details for product ID: %s", order.ProductID)
		product, err := c.productClient.GetProduct(order.ProductID)
		if err != nil {
			c.logger.Error("Failed to get product details for product ID: %s: %v. Message will be requeued", order.ProductID, err)
			msg.Nack(false, true)
			continue
		}

		c.logger.Debug("Received product details - Name: %s, Price: %.2f, Stock: %d",
			product.Product.Name, product.Product.Price, product.Product.Stock)

		order.TotalPrice = product.Product.Price * float64(order.Quantity)
		c.logger.Info("Calculated total price: %.2f for quantity: %d", order.TotalPrice, order.Quantity)

		c.logger.Debug("Updating order total price in database for order ID: %s", order.ID.Hex())
		err = c.repo.UpdateTotalPrice(order.ID.Hex(), order.TotalPrice)
		if err != nil {
			c.logger.Error("Failed to update total price in DB for order ID: %s: %v. Message will be rejected", order.ID.Hex(), err)
			msg.Nack(false, false)
			continue
		}
		c.logger.Info("Successfully updated order total price in database")

		c.logger.Debug("Attempting to decrease stock by %d for product ID: %s", order.Quantity, order.ProductID)
		err = c.productClient.DecreaseStock(order.ProductID, int32(order.Quantity))
		if err != nil {
			c.logger.Error("Failed to decrease stock for product ID: %s: %v. Message will be REJECTED", order.ProductID, err)
			msg.Nack(false, false)
			continue
		}
		c.logger.Info("Successfully decreased stock for product ID: %s", order.ProductID)

		err = msg.Ack(false)
		if err != nil {
			c.logger.Error("Failed to acknowledge message for order ID: %s: %v", order.ID.Hex(), err)
			continue
		}
		// After successful order processing
		c.logger.Info("Order %s successfully processed for user %s (%s). Product: %s, Quantity: %d, Total: %.2f",
			order.ID.Hex(), user.Id, user.Email, product.Product.Name, order.Quantity, order.TotalPrice)
		c.logger.Debug("User details - ID: %s, Email: %s", user.Id, user.Email)
		c.logger.Info("Successfully processed order ID: %s. Product: %s, Quantity: %d, Total: %.2f",
			order.ID.Hex(), product.Product.Name, order.Quantity, order.TotalPrice)
	}

	return nil
}

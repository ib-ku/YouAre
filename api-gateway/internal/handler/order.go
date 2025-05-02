package handler

import (
	"context"
	"net/http"
	orderpb "order_service/pkg/gen/order"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderHandler struct {
	orderClient orderpb.OrderServiceClient
}

func NewOrderHandler(orderClient orderpb.OrderServiceClient) *OrderHandler {
	return &OrderHandler{orderClient: orderClient}
}

// CreateOrder creates a new order
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req orderpb.CreateOrderRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Call Order Service
	resp, err := h.orderClient.CreateOrder(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetOrdersByUser fetches all orders for a user
func (h *OrderHandler) GetOrdersByUser(c *gin.Context) {
	userID := c.Param("user_id")

	// Call Order Service to get all orders
	resp, err := h.orderClient.GetAllOrders(context.Background(), &emptypb.Empty{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Filter orders by user_id
	filtered := []*orderpb.Order{}
	for _, order := range resp.Orders {
		if order.UserId == userID {
			filtered = append(filtered, order)
		}
	}

	c.JSON(http.StatusOK, gin.H{"orders": filtered})
}

// DeleteOrder deletes an order by ID
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	orderID := c.Param("order_id")

	// Corrected: Use OrderRequest instead of DeleteOrderRequest
	req := &orderpb.OrderRequest{Id: orderID}

	// Call Order Service to delete the order
	_, err := h.orderClient.DeleteOrder(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted"})
}

package handler

import (
	"net/http"

	productpb "product-service/pkg/gen/product"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductHandler struct {
	client productpb.ProductServiceClient
}

func NewProductHandler(client productpb.ProductServiceClient) *ProductHandler {
	return &ProductHandler{client: client}
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.client.Get(c, &productpb.ProductRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get product"})
		return
	}
	c.JSON(http.StatusOK, resp.Product)
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	resp, err := h.client.GetAll(c, &emptypb.Empty{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list products"})
		return
	}
	c.JSON(http.StatusOK, resp.Products)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req productpb.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	resp, err := h.client.Create(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
		return
	}
	c.JSON(http.StatusCreated, resp.Product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	_, err := h.client.Delete(c, &productpb.ProductRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
		return
	}
	c.Status(http.StatusNoContent)
}

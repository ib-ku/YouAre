package grpc

import (
	"context"
	productpb "product-service/pkg/gen/product"
)

type ProductClient struct {
	client productpb.ProductServiceClient
}

func NewProductClient(client productpb.ProductServiceClient) *ProductClient {
	return &ProductClient{
		client: client,
	}
}

func (p *ProductClient) DecreaseStock(productID string, quantity int32) error {
	_, err := p.client.Decrease(context.TODO(), &productpb.DecreaseRequest{
		Id:       productID,
		Quantity: quantity,
	})
	return err
}

func (p *ProductClient) GetProduct(productID string) (*productpb.ProductResponse, error) {
	return p.client.Get(context.TODO(), &productpb.ProductRequest{Id: productID})
}

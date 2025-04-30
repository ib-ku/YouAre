package transport

import (
	"YouAre/product_service/application"
	"YouAre/product_service/domain"
	"YouAre/product_service/proto"
	"context"
)

type ProductServer struct {
	service *application.ProductService
}

func NewProductServer(service *application.ProductService) *ProductServer {
	return &ProductServer{service: service}
}

func (s *ProductServer) CreateProduct(ctx context.Context, req *proto.Product) (*proto.ProductResponse, error) {
	product := &domain.Product{
		ID:    req.Id,
		Name:  req.Name,
		Price: req.Price,
	}
	err := s.service.Create(product)
	if err != nil {
		return nil, err
	}
	return &proto.ProductResponse{Success: true, Message: "Product created successfully"}, nil
}

func (s *ProductServer) GetProduct(ctx context.Context, req *proto.ProductRequest) (*proto.Product, error) {
	product, err := s.service.Get(req.Id)
	if err != nil {
		return nil, err
	}
	return &proto.Product{
		Id:    product.ID,
		Name:  product.Name,
		Price: product.Price,
	}, nil
}

func (s *ProductServer) GetAllProducts(ctx context.Context, req *proto.Empty) (*proto.Products, error) {
	products, _ := s.service.GetAll()
	var protoProducts []*proto.Product
	for _, p := range products {
		protoProducts = append(protoProducts, &proto.Product{
			Id:    p.ID,
			Name:  p.Name,
			Price: p.Price,
		})
	}
	return &proto.Products{Product: protoProducts}, nil
}

func (s *ProductServer) UpdateProduct(ctx context.Context, req *proto.Product) (*proto.ProductResponse, error) {
	product := &domain.Product{
		ID:    req.Id,
		Name:  req.Name,
		Price: req.Price,
	}
	err := s.service.Update(product)
	if err != nil {
		return nil, err
	}
	return &proto.ProductResponse{Success: true, Message: "Product updated successfully"}, nil
}

func (s *ProductServer) DeleteProduct(ctx context.Context, req *proto.ProductRequest) (*proto.ProductResponse, error) {
	err := s.service.Delete(req.Id)
	if err != nil {
		return nil, err
	}
	return &proto.ProductResponse{Success: true, Message: "Product deleted successfully"}, nil
}

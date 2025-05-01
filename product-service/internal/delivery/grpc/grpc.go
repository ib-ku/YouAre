package grpc

import (
	"context"
	"product-service/internal/entity"
	"product-service/internal/usecase"
	productpb "product-service/pkg/gen/product"

	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type ProductHandler struct {
	productpb.UnimplementedProductServiceServer
	usecase usecase.ProductUseCase
}

func NewProductHandler(productUC usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		usecase: productUC,
	}
}

func (h *ProductHandler) Create(ctx context.Context, req *productpb.CreateRequest) (*productpb.ProductResponse, error) {
	newProduct := &entity.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: int(req.Stock),
	}

	_, err := h.usecase.CreateProduct(newProduct.Name, newProduct.Price, newProduct.Stock)
	if err != nil {
		return nil, err
	}

	return &productpb.ProductResponse{
		Product: &productpb.Product{
			Name:  newProduct.Name,
			Price: newProduct.Price,
			Stock: int32(newProduct.Stock),
		},
	}, nil
}

func (h *ProductHandler) Get(ctx context.Context, req *productpb.ProductRequest) (*productpb.ProductResponse, error) {
	product, err := h.usecase.GetProduct(req.Id)
	if err != nil {
		return nil, err
	}

	return &productpb.ProductResponse{
		Product: &productpb.Product{
			Id:    product.ID.Hex(),
			Name:  product.Name,
			Price: product.Price,
			Stock: int32(product.Stock),
		},
	}, nil
}
func (h *ProductHandler) GetAll(ctx context.Context, _ *emptypb.Empty) (*productpb.ProductListResponse, error) {
	products, err := h.usecase.GetAllProducts()
	if err != nil {
		return nil, err
	}

	var pbProducts []*productpb.Product
	for _, p := range products {
		pbProducts = append(pbProducts, &productpb.Product{
			Id:    p.ID.Hex(),
			Name:  p.Name,
			Price: p.Price,
			Stock: int32(p.Stock),
		})
	}
	return &productpb.ProductListResponse{Products: pbProducts}, nil
}

func (h *ProductHandler) Update(ctx context.Context, req *productpb.UpdateRequest) (*productpb.ProductResponse, error) {
	updated, err := h.usecase.UpdateProduct(req.Id, req.Name, req.Price, int(req.Stock))
	if err != nil {
		return nil, err
	}

	return &productpb.ProductResponse{
		Product: &productpb.Product{
			Id:    updated.ID.Hex(),
			Name:  updated.Name,
			Price: updated.Price,
			Stock: int32(updated.Stock),
		},
	}, nil
}

func (h *ProductHandler) Decrease(ctx context.Context, req *productpb.DecreaseRequest) (*productpb.ProductResponse, error) {
	p, err := h.usecase.DecreaseStock(req.Id, int(req.Quantity))
	if err != nil {
		return nil, err
	}

	return &productpb.ProductResponse{
		Product: &productpb.Product{
			Id:    p.ID.Hex(),
			Name:  p.Name,
			Price: p.Price,
			Stock: int32(p.Stock),
		},
	}, nil
}

func (h *ProductHandler) Delete(ctx context.Context, req *productpb.ProductRequest) (*emptypb.Empty, error) {
	err := h.usecase.DeleteProduct(req.Id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

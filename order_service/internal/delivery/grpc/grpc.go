package grpc

import (
	"context"
	"order_service/internal/entity"
	"order_service/internal/usecase"
	orderpb "order_service/pkg/gen/order"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	orderpb.UnimplementedOrderServiceServer
	uc *usecase.OrderUseCase
}

func NewHandler(uc *usecase.OrderUseCase) orderpb.OrderServiceServer {
	return &Handler{uc: uc}
}
func (h *Handler) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.OrderResponse, error) {
	newOrder := &entity.Order{
		UserID:    req.UserId,
		ProductID: req.ProductId,
		Quantity:  int(req.Quantity),
		// total_price and created_at will be set in usecase or entity logic
	}

	createdOrder, err := h.uc.CreateOrder(newOrder)
	if err != nil {
		return nil, err
	}

	return &orderpb.OrderResponse{
		Order: &orderpb.Order{
			Id:         createdOrder.ID.Hex(),
			UserId:     createdOrder.UserID,
			ProductId:  createdOrder.ProductID,
			Quantity:   int32(createdOrder.Quantity),
			TotalPrice: createdOrder.TotalPrice,
			CreatedAt:  timestamppb.New(createdOrder.CreatedAt),
		},
	}, nil
}

func (h *Handler) GetOrder(ctx context.Context, req *orderpb.OrderRequest) (*orderpb.OrderResponse, error) {
	order, err := h.uc.GetOrder(req.Id)
	if err != nil {
		return nil, err
	}

	return &orderpb.OrderResponse{
		Order: &orderpb.Order{
			Id:         order.ID.Hex(),
			UserId:     order.UserID,
			ProductId:  order.ProductID,
			Quantity:   int32(order.Quantity),
			TotalPrice: order.TotalPrice,
			CreatedAt:  timestamppb.New(order.CreatedAt),
		},
	}, nil
}

func (h *Handler) GetAllOrders(ctx context.Context, _ *emptypb.Empty) (*orderpb.OrderListResponse, error) {
	orders, err := h.uc.GetAllOrders()
	if err != nil {
		return nil, err
	}

	var pbOrders []*orderpb.Order
	for _, o := range orders {
		pbOrders = append(pbOrders, &orderpb.Order{
			Id:         o.ID.Hex(),
			UserId:     o.UserID,
			ProductId:  o.ProductID,
			Quantity:   int32(o.Quantity),
			TotalPrice: o.TotalPrice,
			CreatedAt:  timestamppb.New(o.CreatedAt),
		})
	}

	return &orderpb.OrderListResponse{Orders: pbOrders}, nil
}

func (h *Handler) DeleteOrder(ctx context.Context, req *orderpb.OrderRequest) (*emptypb.Empty, error) {
	err := h.uc.DeleteOrder(req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

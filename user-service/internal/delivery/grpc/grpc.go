package grpc

import (
	"context"
	"user-service/internal/entity"
	"user-service/internal/usecase"
	userpb "user-service/pkg/gen/user"

	"google.golang.org/protobuf/types/known/emptypb"
)

type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	usecase     usecase.UserUsecase
	authUsecase usecase.AuthUsecase
}

func NewUserHandler(userUC usecase.UserUsecase, authUC usecase.AuthUsecase) *UserHandler {
	return &UserHandler{
		usecase:     userUC,
		authUsecase: authUC,
	}
}

func (h *UserHandler) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.UserResponse, error) {
	newUser := &entity.User{
		Email:    req.Email,
		Password: req.Password,
	}

	_, err := h.authUsecase.Register(newUser.Email, newUser.Password)
	if err != nil {
		return nil, err
	}

	return &userpb.UserResponse{
		Id:    newUser.ID.Hex(),
		Email: newUser.Email,
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	tokenPair, err := h.authUsecase.Login(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	user, err := h.usecase.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	return &userpb.LoginResponse{
		AccessToken: tokenPair.AccessToken,
		User: &userpb.UserResponse{
			Id:    user.ID.Hex(),
			Email: user.Email,
		},
	}, nil
}

func (h *UserHandler) GetProfile(ctx context.Context, req *userpb.ProfileRequest) (*userpb.UserResponse, error) {
	user, err := h.usecase.GetUserById(req.UserId)
	if err != nil {
		return nil, err
	}

	return &userpb.UserResponse{
		Id:    user.ID.Hex(),
		Email: user.Email,
	}, nil
}

func (h *UserHandler) GetAllProfiles(ctx context.Context, req *emptypb.Empty) (*userpb.UserListResponse, error) {
	users, err := h.usecase.GetAll()
	if err != nil {
		return nil, err
	}

	var userResponses []*userpb.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, &userpb.UserResponse{
			Id:    user.ID.Hex(),
			Email: user.Email,
		})
	}

	return &userpb.UserListResponse{
		Users: userResponses,
	}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *userpb.ProfileRequest) (*emptypb.Empty, error) {
	err := h.usecase.DeleteUser(req.UserId)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	pbUser "github.com/relaunch-cot/lib-relaunch-cot/proto/user"
	"github.com/relaunch-cot/service-user/handler"
)

type userResource struct {
	handler *handler.Handlers
	pbUser.UnimplementedUserServiceServer
}

func (r *userResource) CreateUser(ctx context.Context, in *pbUser.CreateUserRequest) (*empty.Empty, error) {
	err := r.handler.User.CreateUser(&ctx, in.Name, in.Email, in.Password)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (r *userResource) LoginUser(ctx context.Context, in *pbUser.LoginUserRequest) (*pbUser.LoginUserResponse, error) {
	loginUserResponse, err := r.handler.User.LoginUser(&ctx, in.Email, in.Password)
	if err != nil {
		return nil, err
	}

	return loginUserResponse, nil
}

func (r *userResource) UpdateUserPassword(ctx context.Context, in *pbUser.UpdateUserPasswordRequest) (*empty.Empty, error) {
	err := r.handler.User.UpdateUserPassword(&ctx, in)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (r *userResource) UpdateUser(ctx context.Context, in *pbUser.UpdateUserRequest) (*empty.Empty, error) {
	err := r.handler.User.UpdateUser(&ctx, in)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (r *userResource) DeleteUser(ctx context.Context, in *pbUser.DeleteUserRequest) (*empty.Empty, error) {
	err := r.handler.User.DeleteUser(&ctx, in)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func NewUserServer(handler *handler.Handlers) pbUser.UserServiceServer {
	return &userResource{
		handler: handler,
	}
}

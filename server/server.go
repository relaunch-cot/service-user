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

func NewUserServer(handler *handler.Handlers) pbUser.UserServiceServer {
	return &userResource{
		handler: handler,
	}
}

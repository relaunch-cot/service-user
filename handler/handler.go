package handler

import (
	"context"
	"github.com/relaunch-cot/bff-relaunch/grpc"
	pb "github.com/relaunch-cot/lib-relaunch-cot/proto/user"
	"github.com/relaunch-cot/service-user/repositories"
)

type IUserHandler interface {
	CreateUser(ctx *context.Context, name, email, password string) error
	LoginUser(ctx *context.Context, email, password string) (pb.LoginUserResponse, error)
	UpdateUserPassword(ctx *context.Context, in *pb.UpdateUserPasswordRequest) error
}

type resource struct {
	repositories *repositories.Repositories
	grpc         grpc.Grpc
}

func (r *resource) CreateUser(ctx *context.Context, name, email, password string) error {
	err := r.repositories.Mysql.CreateUser(ctx, name, email, password)
	if err != nil {
		return err
	}

	return nil
}

func (r *resource) LoginUser(ctx *context.Context, email, password string) (pb.LoginUserResponse, error) {
	loginUserResponse, err := r.repositories.Mysql.LoginUser(ctx, email, password)
	if err != nil {
		return pb.LoginUserResponse{}, err
	}

	return loginUserResponse, nil
}

func (r *resource) UpdateUserPassword(ctx *context.Context, in *pb.UpdateUserPasswordRequest) error {
	err := r.repositories.Mysql.UpdateUserPassword(ctx, in.Email, in.CurrentPassword, in.NewPassword)
	if err != nil {
		return err
	}

	return nil
}

func NewUserHandler(repositories *repositories.Repositories) IUserHandler {
	return &resource{
		repositories: repositories,
	}
}

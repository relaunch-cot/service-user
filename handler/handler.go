package handler

import (
	"context"
	"github.com/relaunch-cot/bff/grpc"
	"github.com/relaunch-cot/service-user/repositories"
)

type IUserHandler interface {
	CreateUser(ctx *context.Context, name, email, password string) error
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

func NewUserHandler(repositories *repositories.Repositories) IUserHandler {
	return &resource{
		repositories: repositories,
	}
}

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
	err := r.handler.User.CreateUser(&ctx, in.Name, in.Email, in.Password, in.Type, in.Settings)
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

func (r *userResource) GenerateReportFromJSON(ctx context.Context, in *pbUser.GenerateReportRequest) (*pbUser.GenerateReportResponse, error) {
	pdfBytes, err := r.handler.User.GenerateReportFromJSON(&ctx, in.JsonData)
	if err != nil {
		return nil, err
	}

	return &pbUser.GenerateReportResponse{
		PdfData: pdfBytes,
	}, nil
}

func (r *userResource) SendPasswordRecoveryEmail(ctx context.Context, in *pbUser.SendPasswordRecoveryEmailRequest) (*empty.Empty, error) {
	err := r.handler.User.SendPasswordRecoveryEmail(&ctx, in.Email, in.RecoveryLink)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (r *userResource) GetUserProfile(ctx context.Context, in *pbUser.GetUserProfileRequest) (*pbUser.GetUserProfileResponse, error) {
	response, err := r.handler.User.GetUserProfile(&ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *userResource) GetUserType(ctx context.Context, in *pbUser.GetUserTypeRequest) (*pbUser.GetUserTypeResponse, error) {
	response, err := r.handler.User.GetUserType(&ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewUserServer(handler *handler.Handlers) pbUser.UserServiceServer {
	return &userResource{
		handler: handler,
	}
}

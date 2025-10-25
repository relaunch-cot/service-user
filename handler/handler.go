package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	libModels "github.com/relaunch-cot/lib-relaunch-cot/models"
	"github.com/relaunch-cot/lib-relaunch-cot/proto/base_models"
	pb "github.com/relaunch-cot/lib-relaunch-cot/proto/user"
	"github.com/relaunch-cot/service-user/config"
	"github.com/relaunch-cot/service-user/repositories"
	"github.com/relaunch-cot/service-user/resource/transformer"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IUserHandler interface {
	CreateUser(ctx *context.Context, name, email, password, userType string, settings *base_models.UserSettings) error
	LoginUser(ctx *context.Context, email, password string) (*pb.LoginUserResponse, error)
	UpdateUser(ctx *context.Context, in *pb.UpdateUserRequest) error
	UpdateUserPassword(ctx *context.Context, in *pb.UpdateUserPasswordRequest) error
	DeleteUser(ctx *context.Context, in *pb.DeleteUserRequest) error
	GenerateReportFromJSON(ctx *context.Context, jsonData string) ([]byte, error)
	SendPasswordRecoveryEmail(ctx *context.Context, email, recoveryLink string) error
	GetUserProfile(ctx *context.Context, userId string) (*pb.GetUserProfileResponse, error)
	GetUserType(ctx *context.Context, userId string) (*pb.GetUserTypeResponse, error)
}

type resource struct {
	repositories *repositories.Repositories
}

func (r *resource) CreateUser(ctx *context.Context, name, email, password, userType string, settings *base_models.UserSettings) error {
	userId := uuid.New()

	err := r.repositories.Mysql.CreateUser(ctx, userId.String(), name, email, password, userType, settings)
	if err != nil {
		return err
	}

	return nil
}

func (r *resource) LoginUser(ctx *context.Context, email, password string) (*pb.LoginUserResponse, error) {
	user, err := r.repositories.Mysql.LoginUser(ctx, email, password)
	if err != nil {
		return nil, err
	}

	tokenString, err := createToken(user.UserId)
	if err != nil {
		return nil, err
	}

	loginUserResponse := &pb.LoginUserResponse{
		Token: tokenString,
	}

	return loginUserResponse, nil
}

func (r *resource) UpdateUserPassword(ctx *context.Context, in *pb.UpdateUserPasswordRequest) error {
	err := r.repositories.Mysql.UpdateUserPassword(ctx, in.UserId, in.NewPassword)
	if err != nil {
		return err
	}

	return nil
}

func (r *resource) UpdateUser(ctx *context.Context, in *pb.UpdateUserRequest) error {
	err := r.repositories.Mysql.UpdateUser(ctx, in.UserId, in.NewUser)
	if err != nil {
		return err
	}

	return nil
}

func (r *resource) DeleteUser(ctx *context.Context, in *pb.DeleteUserRequest) error {
	err := r.repositories.Mysql.DeleteUser(ctx, in.Email, in.Password)
	if err != nil {
		return err
	}

	return nil
}

func (r *resource) GenerateReportFromJSON(ctx *context.Context, jsonData string) ([]byte, error) {
	var reportData libModels.ReportData
	err := json.Unmarshal([]byte(jsonData), &reportData)
	if err != nil {
		return nil, status.Error(codes.Internal, "error unmarshalling report data. Details: "+err.Error())
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Título
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, reportData.Title)
	pdf.Ln(15)

	if reportData.Subtitle != "" {
		pdf.SetFont("Arial", "", 12)
		pdf.Cell(190, 8, reportData.Subtitle)
		pdf.Ln(10)
	}

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(190, 8, fmt.Sprintf("Gerado em: %s", time.Now().Format("02/01/2006 15:04:05")))
	pdf.Ln(15)

	if len(reportData.Headers) > 0 {
		pdf.SetFont("Arial", "B", 10)
		cellWidth := 190.0 / float64(len(reportData.Headers))

		for _, header := range reportData.Headers {
			pdf.Cell(cellWidth, 8, header)
		}
		pdf.Ln(10)

		pdf.SetFont("Arial", "", 9)
		for _, row := range reportData.Rows {
			for i, cell := range row {
				if i < len(reportData.Headers) {
					pdf.Cell(cellWidth, 7, cell)
				}
			}
			pdf.Ln(7)
		}
	}

	if reportData.Footer != "" {
		pdf.Ln(10)
		pdf.SetFont("Arial", "I", 10)
		pdf.Cell(190, 8, reportData.Footer)
	}

	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(190, 8, fmt.Sprintf("Total de registros: %d", len(reportData.Rows)))

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, status.Error(codes.Internal, "error generating report. Details: "+err.Error())
	}

	return buf.Bytes(), nil
}

func (r *resource) SendPasswordRecoveryEmail(ctx *context.Context, email, recoveryLink string) error {
	name, err := r.repositories.Mysql.SendPasswordRecoveryEmail(ctx, email)
	if err != nil {
		return err
	}
	from := mail.NewEmail(config.NAME, config.EMAIL)
	subject := "Recuperação de Senha"
	to := mail.NewEmail(*name, email)

	plainTextContent := fmt.Sprintf("Olá %s,\n\nClique no link abaixo para redefinir sua senha:\n%s\n\nSe você não solicitou, ignore este e-mail.", *name, recoveryLink)

	htmlContent := fmt.Sprintf(`
        <p>Olá %s,</p>
        <p>Clique no link abaixo para redefinir sua senha:</p>
        <p><a href="%s">Recuperar Senha</a></p>
        <p>Se você não solicitou, ignore este e-mail.</p>
    `, *name, recoveryLink)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	replyTo := mail.NewEmail("ReLaunch Support", "support@relaunch.com.br")
	message.SetReplyTo(replyTo)

	client := sendgrid.NewSendClient(config.SENDGRID_API_KEY)
	response, err := client.Send(message)
	if err != nil {
		return status.Error(codes.Internal, "error sending email. Details: "+err.Error())
	}

	if response.StatusCode >= 400 {
		return status.Error(codes.Code(response.StatusCode), "error sending email. Details: "+response.Body)
	}

	return nil
}

func (r *resource) GetUserProfile(ctx *context.Context, userId string) (*pb.GetUserProfileResponse, error) {
	mysqlResponse, err := r.repositories.Mysql.GetUserProfile(ctx, userId)
	if err != nil {
		return nil, err
	}

	baseModelsUser, err := transformer.GetUserProfileToBaseModels(mysqlResponse)
	if err != nil {
		return nil, err
	}

	getUserProfileResponse := &pb.GetUserProfileResponse{
		User: baseModelsUser,
	}

	return getUserProfileResponse, nil
}

func (r *resource) GetUserType(ctx *context.Context, userId string) (*pb.GetUserTypeResponse, error) {
	mysqlResponse, err := r.repositories.Mysql.GetUserType(ctx, userId)
	if err != nil {
		return nil, err
	}

	getUserTypeResponse := &pb.GetUserTypeResponse{
		UserType: *mysqlResponse,
	}

	return getUserTypeResponse, nil
}

var secretKey = []byte(config.JWT_SECRET)

func createToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": userId,
			"exp":    time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", status.Error(codes.Internal, "error signing token. Details: "+err.Error())
	}

	tokenString = fmt.Sprintf(`Bearer %s`, tokenString)

	return tokenString, nil
}

func NewUserHandler(repositories *repositories.Repositories) IUserHandler {
	return &resource{
		repositories: repositories,
	}
}

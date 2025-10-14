package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
	pb "github.com/relaunch-cot/lib-relaunch-cot/proto/user"
	"github.com/relaunch-cot/service-user/config"
	"github.com/relaunch-cot/service-user/repositories"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type ReportData struct {
	Title    string     `json:"title"`
	Subtitle string     `json:"subtitle,omitempty"`
	Headers  []string   `json:"headers"`
	Rows     [][]string `json:"rows"`
	Footer   string     `json:"footer,omitempty"`
}

type IUserHandler interface {
	CreateUser(ctx *context.Context, name, email, password string) error
	LoginUser(ctx *context.Context, email, password string) (*pb.LoginUserResponse, error)
	UpdateUser(ctx *context.Context, in *pb.UpdateUserRequest) error
	UpdateUserPassword(ctx *context.Context, in *pb.UpdateUserPasswordRequest) error
	DeleteUser(ctx *context.Context, in *pb.DeleteUserRequest) error
	GenerateReportFromJSON(ctx *context.Context, jsonData string) ([]byte, error)
	SendPasswordRecoveryEmail(ctx *context.Context, email, recoveryLink string) error
	CreateNewChat(ctx *context.Context, createdBy int64, userIds []int64) error
}

type resource struct {
	repositories *repositories.Repositories
}

func (r *resource) CreateUser(ctx *context.Context, name, email, password string) error {
	err := r.repositories.Mysql.CreateUser(ctx, name, email, password)
	if err != nil {
		return err
	}

	return nil
}

func (r *resource) LoginUser(ctx *context.Context, email, password string) (*pb.LoginUserResponse, error) {
	loginUserResponse, err := r.repositories.Mysql.LoginUser(ctx, email, password)
	if err != nil {
		return nil, err
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
	err := r.repositories.Mysql.UpdateUser(ctx, in.Password, in.UserId, in.NewUser)
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
	var reportData ReportData
	err := json.Unmarshal([]byte(jsonData), &reportData)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer parse do JSON: %v", err)
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
		return nil, err
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
		return fmt.Errorf("erro ao enviar email: %w", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("erro no envio do email: %s", response.Body)
	}

	return nil
}

func (r *resource) CreateNewChat(ctx *context.Context, createdBy int64, userIds []int64) error {
	err := r.repositories.Mysql.CreateNewChat(ctx, createdBy, userIds)
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

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
	pb "github.com/relaunch-cot/lib-relaunch-cot/proto/user"
	"github.com/relaunch-cot/service-user/repositories"
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
	err := r.repositories.Mysql.UpdateUserPassword(ctx, in.Email, in.CurrentPassword, in.NewPassword)
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
	// Parse do JSON recebido
	var reportData ReportData
	err := json.Unmarshal([]byte(jsonData), &reportData)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer parse do JSON: %v", err)
	}

	// Criar PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Título
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, reportData.Title)
	pdf.Ln(15)

	// Subtítulo (se existir)
	if reportData.Subtitle != "" {
		pdf.SetFont("Arial", "", 12)
		pdf.Cell(190, 8, reportData.Subtitle)
		pdf.Ln(10)
	}

	// Data de geração
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(190, 8, fmt.Sprintf("Gerado em: %s", time.Now().Format("02/01/2006 15:04:05")))
	pdf.Ln(15)

	// Cabeçalhos da tabela
	if len(reportData.Headers) > 0 {
		pdf.SetFont("Arial", "B", 10)
		cellWidth := 190.0 / float64(len(reportData.Headers))

		for _, header := range reportData.Headers {
			pdf.Cell(cellWidth, 8, header)
		}
		pdf.Ln(10)

		// Dados das linhas
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

	// Footer (se existir)
	if reportData.Footer != "" {
		pdf.Ln(10)
		pdf.SetFont("Arial", "I", 10)
		pdf.Cell(190, 8, reportData.Footer)
	}

	// Adicionar estatísticas
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(190, 8, fmt.Sprintf("Total de registros: %d", len(reportData.Rows)))

	// Converter para bytes usando buffer
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func NewUserHandler(repositories *repositories.Repositories) IUserHandler {
	return &resource{
		repositories: repositories,
	}
}

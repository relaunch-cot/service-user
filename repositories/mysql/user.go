package mysql

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/relaunch-cot/lib-relaunch-cot/proto/user"
	"github.com/relaunch-cot/lib-relaunch-cot/repositories/mysql"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type mysqlResource struct {
	client *mysql.Client
}

type IMySqlUser interface {
	CreateUser(ctx *context.Context, name, email, password string) error
	LoginUser(ctx *context.Context, email, password string) (pb.LoginUserResponse, error)
}

func (r *mysqlResource) CreateUser(ctx *context.Context, name, email, password string) error {
	basequery := fmt.Sprintf(
		"INSERT INTO user (name, email, password) VALUES('%s', '%s', '%s')",
		name,
		email,
		password,
	)
	rows, err := mysql.DB.QueryContext(*ctx, basequery)
	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}

func (r *mysqlResource) LoginUser(ctx *context.Context, email, password string) (pb.LoginUserResponse, error) {
	basequery := fmt.Sprintf(`SELECT * FROM user WHERE email = '%s' AND password = '%s'`, email, password)
	rows, err := mysql.DB.QueryContext(*ctx, basequery)
	if err != nil {
		return pb.LoginUserResponse{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return pb.LoginUserResponse{}, errors.New("user not found")
	}

	tokenString, err := createToken(email)
	if err != nil {
		return pb.LoginUserResponse{}, err
	}

	loginUserResponse := pb.LoginUserResponse{
		Token: tokenString,
	}

	return loginUserResponse, nil
}

var secretKey = []byte("secret-key")

func createToken(userEmail string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userEmail": userEmail,
			"exp":       time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	tokenString = fmt.Sprintf(`Bearer: %s`, tokenString)

	return tokenString, nil
}

func NewMysqlRepository(client *mysql.Client) IMySqlUser {
	return &mysqlResource{
		client: client,
	}
}

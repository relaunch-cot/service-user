package mysql

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/relaunch-cot/lib-relaunch-cot/proto/user"
	"github.com/relaunch-cot/lib-relaunch-cot/repositories/mysql"
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type mysqlResource struct {
	client *mysql.Client
}

type IMySqlUser interface {
	CreateUser(ctx *context.Context, name, email, password string) error
	LoginUser(ctx *context.Context, email, password string) (pb.LoginUserResponse, error)
	UpdateUserPassword(ctx *context.Context, email, currentPassword, newPassword string) error
}

func (r *mysqlResource) CreateUser(ctx *context.Context, name, email, password string) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	basequery := fmt.Sprintf(
		"INSERT INTO user (name, email, password) VALUES('%s', '%s', '%s')",
		name,
		email,
		hashPassword,
	)
	rows, err := mysql.DB.QueryContext(*ctx, basequery)
	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}

type User struct {
	userId         int
	name           string
	email          string
	hashedPassword string
}

func (r *mysqlResource) LoginUser(ctx *context.Context, email, password string) (pb.LoginUserResponse, error) {
	var user User

	basequery := fmt.Sprintf(`SELECT * FROM user WHERE email = '%s'`, email)
	rows, err := mysql.DB.QueryContext(*ctx, basequery)
	if err != nil {
		return pb.LoginUserResponse{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return pb.LoginUserResponse{}, errors.New("user not found")
	}

	err = rows.Scan(&user.userId, &user.name, &user.hashedPassword, &user.email)
	if err != nil {
		return pb.LoginUserResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.hashedPassword), []byte(password))
	if err != nil {
		return pb.LoginUserResponse{}, errors.New("wrong password")
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

func (r *mysqlResource) UpdateUserPassword(ctx *context.Context, email, currentPassword, newPassword string) error {
	var user User

	queryValidateUser := fmt.Sprintf(`SELECT * FROM user WHERE email = '%s'`, email)
	rows, err := mysql.DB.QueryContext(*ctx, queryValidateUser)
	if err != nil {
		return err
	}

	defer rows.Close()

	if !rows.Next() {
		return errors.New("user not found")
	}

	err = rows.Scan(&user.userId, &user.name, &user.hashedPassword, &user.email)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.hashedPassword), []byte(currentPassword))
	if err != nil {
		return errors.New("wrong password")
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 14)
	if err != nil {
		return err
	}

	updateQuery := fmt.Sprintf(`UPDATE user SET password = '%s' WHERE email = '%s'`, newHashedPassword, email)
	_, err = mysql.DB.ExecContext(*ctx, updateQuery)
	if err != nil {
		return err
	}

	return nil
}

func NewMysqlRepository(client *mysql.Client) IMySqlUser {
	return &mysqlResource{
		client: client,
	}
}

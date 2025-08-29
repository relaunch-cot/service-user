package mysql

import (
	"context"
	"errors"
	"fmt"
	"time"

	userModel "github.com/relaunch-cot/lib-relaunch-cot/models/user"
	pb "github.com/relaunch-cot/lib-relaunch-cot/proto/user"
	"github.com/relaunch-cot/lib-relaunch-cot/repositories/mysql"
	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v5"
)

type mysqlResource struct {
	client *mysql.Client
}

type IMySqlUser interface {
	CreateUser(ctx *context.Context, name, email, password string) error
	LoginUser(ctx *context.Context, email, password string) (*pb.LoginUserResponse, error)
	UpdateUser(ctx *context.Context, currentUser, newUser *pb.User) error
	UpdateUserPassword(ctx *context.Context, email, currentPassword, newPassword string) error
}

func (r *mysqlResource) CreateUser(ctx *context.Context, name, email, password string) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	basequery := fmt.Sprintf(
		"INSERT INTO users (name, email, password) VALUES('%s', '%s', '%s')",
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

func (r *mysqlResource) LoginUser(ctx *context.Context, email, password string) (*pb.LoginUserResponse, error) {
	var User userModel.User

	basequery := fmt.Sprintf(`SELECT * FROM users WHERE email = '%s'`, email)
	rows, err := mysql.DB.QueryContext(*ctx, basequery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, errors.New("user not found")
	}

	err = rows.Scan(&User.UserId, &User.Name, &User.HashedPassword, &User.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(User.HashedPassword), []byte(password))
	if err != nil {
		return nil, errors.New("wrong password")
	}

	tokenString, err := createToken(email)
	if err != nil {
		return nil, err
	}

	loginUserResponse := &pb.LoginUserResponse{
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

func (r *mysqlResource) UpdateUser(ctx *context.Context, currentUser, newUser *pb.User) error {
	var User userModel.User

	queryValidateUser := fmt.Sprintf(`SELECT * FROM users WHERE email = '%s'`, currentUser.Email)
	rows, err := mysql.DB.QueryContext(*ctx, queryValidateUser)
	if err != nil {
		return err
	}

	defer rows.Close()

	if !rows.Next() {
		return errors.New("user not found")
	}

	err = rows.Scan(&User.UserId, &User.Name, &User.HashedPassword, &User.Email)
	if err != nil {
		return err
	}

	if currentUser.HashedPassword != "" {
		err = bcrypt.CompareHashAndPassword([]byte(User.HashedPassword), []byte(currentUser.HashedPassword))
		if err != nil {
			return errors.New("wrong password")
		}
	}

	var setParts []string

	// Atualizar nome se fornecido
	if newUser.Name != "" && newUser.Name != User.Name {
		setParts = append(setParts, fmt.Sprintf("name = '%s'", newUser.Name))
	}

	// Atualizar email se fornecido
	if newUser.Email != "" && newUser.Email != User.Email {
		setParts = append(setParts, fmt.Sprintf("email = '%s'", newUser.Email))
	}

	// Atualizar senha se fornecida
	if newUser.HashedPassword != "" {
		newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.HashedPassword), 14)
		if err != nil {
			return err
		}
		setParts = append(setParts, fmt.Sprintf("password = '%s'", string(newHashedPassword)))
	}

	// Se não há campos para atualizar, retornar
	if len(setParts) == 0 {
		return errors.New("no fields to update")
	}

	// Construir e executar query de update
	setClause := setParts[0]
	for i := 1; i < len(setParts); i++ {
		setClause += ", " + setParts[i]
	}

	updateQuery := fmt.Sprintf(`UPDATE users SET %s WHERE email = '%s'`, setClause, currentUser.Email)

	_, err = mysql.DB.ExecContext(*ctx, updateQuery)
	if err != nil {
		return err
	}

	return nil
}

func (r *mysqlResource) UpdateUserPassword(ctx *context.Context, email, currentPassword, newPassword string) error {
	var User userModel.User

	queryValidateUser := fmt.Sprintf(`SELECT * FROM users WHERE email = '%s'`, email)
	rows, err := mysql.DB.QueryContext(*ctx, queryValidateUser)
	if err != nil {
		return err
	}

	defer rows.Close()

	if !rows.Next() {
		return errors.New("user not found")
	}

	err = rows.Scan(&User.UserId, &User.Name, &User.Email, &User.HashedPassword)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(User.HashedPassword), []byte(currentPassword))
	if err != nil {
		return errors.New("wrong password")
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 14)
	if err != nil {
		return err
	}

	updateQuery := fmt.Sprintf(`UPDATE users SET password = '%s' WHERE email = '%s'`, newHashedPassword, email)
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

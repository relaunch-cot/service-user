package mysql

import (
	"context"
	"errors"
	"fmt"
	"time"

	userModel "github.com/relaunch-cot/lib-relaunch-cot/models/user"
	pbBaseModels "github.com/relaunch-cot/lib-relaunch-cot/proto/base_models"
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
	UpdateUser(ctx *context.Context, password string, userId int64, newUser *pbBaseModels.User) error
	UpdateUserPassword(ctx *context.Context, userId int64, newPassword string) error
	DeleteUser(ctx *context.Context, email, password string) error
	SendPasswordRecoveryEmail(ctx *context.Context, email string) (*string, error)
	CreateNewChat(ctx *context.Context, createdBy int64, userIds []int64) error
	SendMessage(ctx *context.Context, chatId, senderId int64, messageContent string) error
}

func (r *mysqlResource) CreateUser(ctx *context.Context, name, email, password string) error {
	queryValidation := fmt.Sprintf(`SELECT * FROM users WHERE email = '%s'`, email)
	rows, err := mysql.DB.QueryContext(*ctx, queryValidation)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return errors.New("already exists an user with this email")
	}

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
	rows, err = mysql.DB.QueryContext(*ctx, basequery)
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

	err = rows.Scan(&User.UserId, &User.Name, &User.Email, &User.HashedPassword)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(User.HashedPassword), []byte(password))
	if err != nil {
		return nil, errors.New("wrong password")
	}

	tokenString, err := createToken(User.UserId)
	if err != nil {
		return nil, err
	}

	loginUserResponse := &pb.LoginUserResponse{
		Token: tokenString,
	}

	return loginUserResponse, nil
}

var secretKey = []byte("secret-key")

func createToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": userId,
			"exp":    time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	tokenString = fmt.Sprintf(`Bearer %s`, tokenString)

	return tokenString, nil
}

func (r *mysqlResource) UpdateUser(ctx *context.Context, password string, userId int64, newUser *pbBaseModels.User) error {
	var User userModel.User

	queryValidateUser := fmt.Sprintf(`SELECT * FROM users WHERE userId = '%d'`, userId)
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
	err = bcrypt.CompareHashAndPassword([]byte(User.HashedPassword), []byte(password))
	if err != nil {
		return errors.New("wrong password")
	}

	var setParts []string

	if newUser.Name != "" && newUser.Name != User.Name {
		setParts = append(setParts, fmt.Sprintf("name = '%s'", newUser.Name))
	}

	if newUser.Email != "" && newUser.Email != User.Email {
		queryValidation := fmt.Sprintf(`SELECT * FROM users WHERE email = '%s'`, newUser.Email)
		rows, err := mysql.DB.QueryContext(*ctx, queryValidation)
		if err != nil {
			return err
		}
		defer rows.Close()

		if rows.Next() {
			return errors.New("already exists an user with this email")
		}
		setParts = append(setParts, fmt.Sprintf("email = '%s'", newUser.Email))
	}

	if len(setParts) == 0 {
		return errors.New("no fields to update")
	}

	setClause := setParts[0]
	for i := 1; i < len(setParts); i++ {
		setClause += ", " + setParts[i]
	}

	updateQuery := fmt.Sprintf(`UPDATE users SET %s WHERE userId = '%d'`, setClause, userId)

	_, err = mysql.DB.ExecContext(*ctx, updateQuery)
	if err != nil {
		return err
	}

	return nil
}

func (r *mysqlResource) UpdateUserPassword(ctx *context.Context, userId int64, newPassword string) error {
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 14)
	if err != nil {
		return err
	}

	updateQuery := fmt.Sprintf(`UPDATE users SET password = '%s' WHERE userId = '%d'`, newHashedPassword, userId)
	_, err = mysql.DB.ExecContext(*ctx, updateQuery)
	if err != nil {
		return err
	}

	return nil
}

func (r *mysqlResource) DeleteUser(ctx *context.Context, email, password string) error {
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

	err = bcrypt.CompareHashAndPassword([]byte(User.HashedPassword), []byte(password))
	if err != nil {
		return errors.New("wrong password")
	}

	deleteQuery := fmt.Sprintf(`DELETE FROM users WHERE userId = '%d'`, User.UserId)
	_, err = mysql.DB.ExecContext(*ctx, deleteQuery)
	if err != nil {
		return err
	}

	return nil
}

func (r *mysqlResource) SendPasswordRecoveryEmail(ctx *context.Context, email string) (*string, error) {
	var User userModel.User

	queryValidateUser := fmt.Sprintf(`SELECT * FROM users WHERE email = '%s' LIMIT 1`, email)
	rows, err := mysql.DB.QueryContext(*ctx, queryValidateUser)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	if !rows.Next() {
		return nil, errors.New("user not found")
	}

	err = rows.Scan(&User.UserId, &User.Name, &User.Email, &User.HashedPassword)
	if err != nil {
		return nil, err
	}

	return &User.Name, nil
}

func (r *mysqlResource) CreateNewChat(ctx *context.Context, createdBy int64, userIds []int64) error {
	currentDate := time.Now()

	queryValidate := fmt.Sprintf(
		`SELECT * 
					FROM chats 
					WHERE user1_id = %d AND user2_id = %d
					OR user1_id = %d AND user2_id = %d`,
		userIds[0], userIds[1], userIds[0], userIds[1],
	)
	rows, err := mysql.DB.QueryContext(*ctx, queryValidate)
	if err != nil {
		return err
	}

	defer rows.Close()
	if rows.Next() {
		return errors.New("already exists an chat with these participants")
	}

	basequery := fmt.Sprintf(
		"INSERT INTO chats (createdBy, user1_id, user2_id, createdAt) VALUES(%d, %d, %d, '%s')",
		createdBy,
		userIds[0],
		userIds[1],
		currentDate.Format("2006-01-02 15:04:05"),
	)

	rows, err = mysql.DB.QueryContext(*ctx, basequery)
	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}

func (r *mysqlResource) SendMessage(ctx *context.Context, chatId, senderId int64, messageContent string) error {
	currentDate := time.Now()

	queryValidation := fmt.Sprintf(
		`SELECT * 
					FROM chats 
					WHERE chatId = %d 
					    AND user1_id = %d OR user2_id = %d`,
		chatId, senderId, senderId,
	)

	rows, err := mysql.DB.QueryContext(*ctx, queryValidation)
	if err != nil {
		return err
	}

	defer rows.Close()
	if !rows.Next() {
		return errors.New("this user is not part of this chat")
	}

	basequery := fmt.Sprintf(
		"INSERT INTO messages (chatId, senderId, content, createdAt) VALUES(%d, %d, '%s', '%s')",
		chatId,
		senderId,
		messageContent,
		currentDate.Format("2006-01-02 15:04:05"),
	)

	rows, err = mysql.DB.QueryContext(*ctx, basequery)
	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}

func NewMysqlRepository(client *mysql.Client) IMySqlUser {
	return &mysqlResource{
		client: client,
	}
}

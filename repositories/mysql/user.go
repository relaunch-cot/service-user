package mysql

import (
	"context"
	"encoding/json"
	"fmt"

	libModels "github.com/relaunch-cot/lib-relaunch-cot/models"
	pbBaseModels "github.com/relaunch-cot/lib-relaunch-cot/proto/base_models"
	"github.com/relaunch-cot/lib-relaunch-cot/repositories/mysql"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mysqlResource struct {
	client *mysql.Client
}

type IMySqlUser interface {
	CreateUser(ctx *context.Context, userId, name, email, password, userType string, settings *pbBaseModels.UserSettings) error
	LoginUser(ctx *context.Context, email, password string) (*libModels.User, error)
	UpdateUser(ctx *context.Context, userId string, newUser *pbBaseModels.User) error
	UpdateUserPassword(ctx *context.Context, userId string, newPassword string) error
	DeleteUser(ctx *context.Context, email, password string) error
	SendPasswordRecoveryEmail(ctx *context.Context, email string) (*string, error)
	GetUserProfile(ctx *context.Context, userId string) (*libModels.User, error)
}

func (r *mysqlResource) CreateUser(ctx *context.Context, userId, name, email, password, userType string, settings *pbBaseModels.UserSettings) error {
	queryValidation := fmt.Sprintf(`SELECT * FROM users WHERE email = '%s'`, email)
	rows, err := mysql.DB.QueryContext(*ctx, queryValidation)
	if err != nil {
		return status.Error(codes.Internal, "error with database. Details: "+err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		return status.Error(codes.AlreadyExists, "already exists an user with this email")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return status.Error(codes.Internal, "error generating password hash. Details: "+err.Error())
	}

	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return status.Error(codes.Internal, "error marshalling settings. Details: "+err.Error())
	}

	basequery := fmt.Sprintf(
		"INSERT INTO users (userId, name, email, password, settings, type) VALUES('%s', '%s', '%s', '%s', '%s', '%s')",
		userId,
		name,
		email,
		hashPassword,
		settingsJSON,
		userType,
	)
	rows, err = mysql.DB.QueryContext(*ctx, basequery)
	if err != nil {
		return status.Error(codes.Internal, "error with database. Details: "+err.Error())
	}

	defer rows.Close()

	return nil
}

func (r *mysqlResource) LoginUser(ctx *context.Context, email, password string) (*libModels.User, error) {
	var User libModels.User

	basequery := fmt.Sprintf(`SELECT u.userId, u.password, u.type FROM users u WHERE u.email = '%s'`, email)
	rows, err := mysql.DB.QueryContext(*ctx, basequery)
	if err != nil {
		return nil, status.Error(codes.Internal, "error with database. Details: "+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	err = rows.Scan(&User.UserId, &User.Password, &User.Type)
	if err != nil {
		return nil, status.Error(codes.Internal, "error scanning mysql row: "+err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(password))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "wrong password")
	}

	return &User, nil
}

func (r *mysqlResource) UpdateUser(ctx *context.Context, userId string, newUser *pbBaseModels.User) error {
	var User libModels.User

	queryValidateUser := fmt.Sprintf(
		`SELECT u.name,
       				   u.email,
       				   u.settings
					FROM users u 
					WHERE u.userId = '%s'`,
		userId,
	)
	rows, err := mysql.DB.QueryContext(*ctx, queryValidateUser)
	if err != nil {
		return status.Error(codes.Internal, "error with database. Details: "+err.Error())
	}

	defer rows.Close()

	if !rows.Next() {
		return status.Error(codes.NotFound, "user not found")
	}

	var settings []byte

	err = rows.Scan(&User.Name, &User.Email, &settings)
	if err != nil {
		return status.Error(codes.Internal, "error scanning mysql row: "+err.Error())
	}
	err = json.Unmarshal(settings, &User.Settings)
	if err != nil {
		return status.Error(codes.Internal, "error unmarshalling settings. Details: "+err.Error())
	}

	var setParts []string

	if newUser.Name != "" && newUser.Name != User.Name {
		setParts = append(setParts, fmt.Sprintf("name = '%s'", newUser.Name))
	}

	if newUser.Email != "" && newUser.Email != User.Email {
		queryValidation := fmt.Sprintf(`SELECT * FROM users WHERE email = '%s'`, newUser.Email)
		rows, err := mysql.DB.QueryContext(*ctx, queryValidation)
		if err != nil {
			return status.Error(codes.Internal, "error with database. Details: "+err.Error())
		}
		defer rows.Close()

		if rows.Next() {
			return status.Error(codes.AlreadyExists, "already exists an user with this email")
		}
		setParts = append(setParts, fmt.Sprintf("email = '%s'", newUser.Email))
	}

	var validateSettings libModels.UserSettings

	if newUser.Settings != nil {
		if newUser.Settings.Phone != "" && newUser.Settings.Phone != User.Settings.Phone {
			validateSettings.Phone = newUser.Settings.Phone
		} else {
			validateSettings.Phone = User.Settings.Phone
		}
		if newUser.Settings.Cpf != "" && newUser.Settings.Cpf != User.Settings.Cpf {
			validateSettings.Cpf = newUser.Settings.Cpf
		} else {
			validateSettings.Cpf = User.Settings.Cpf
		}
		if newUser.Settings.DateOfBirth != "" && newUser.Settings.DateOfBirth != User.Settings.DateOfBirth {
			validateSettings.DateOfBirth = newUser.Settings.DateOfBirth
		} else {
			validateSettings.DateOfBirth = User.Settings.DateOfBirth
		}
		if newUser.Settings.Biography != "" && newUser.Settings.Biography != User.Settings.Biography {
			validateSettings.Biography = newUser.Settings.Biography
		} else {
			validateSettings.Biography = User.Settings.Biography
		}
		if len(newUser.Settings.Skills) != 0 {
			validateSettings.Skills = newUser.Settings.Skills
		} else {
			validateSettings.Skills = User.Settings.Skills
		}

		validateSettingsJSON, err := json.Marshal(validateSettings)
		if err != nil {
			return status.Error(codes.Internal, "error marshalling settings. Details: "+err.Error())
		}
		setParts = append(setParts, fmt.Sprintf("settings = '%s'", validateSettingsJSON))
	}

	if len(setParts) == 0 {
		return status.Error(codes.NotFound, "no fields to update")
	}

	setClause := setParts[0]
	for i := 1; i < len(setParts); i++ {
		setClause += ", " + setParts[i]
	}

	updateQuery := fmt.Sprintf(`UPDATE users SET %s WHERE userId = '%s'`, setClause, userId)

	_, err = mysql.DB.ExecContext(*ctx, updateQuery)
	if err != nil {
		return status.Error(codes.Internal, "error with database. Details: "+err.Error())
	}

	return nil
}

func (r *mysqlResource) UpdateUserPassword(ctx *context.Context, userId, newPassword string) error {
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 14)
	if err != nil {
		return status.Error(codes.Internal, "error generating password hash. Details: "+err.Error())
	}

	updateQuery := fmt.Sprintf(`UPDATE users SET password = '%s' WHERE userId = '%d'`, newHashedPassword, userId)
	_, err = mysql.DB.ExecContext(*ctx, updateQuery)
	if err != nil {
		return status.Error(codes.Internal, "error with database. Details: "+err.Error())
	}

	return nil
}

func (r *mysqlResource) DeleteUser(ctx *context.Context, email, password string) error {
	var User libModels.User

	queryValidateUser := fmt.Sprintf(`SELECT userId, password FROM users WHERE email = '%s'`, email)

	rows, err := mysql.DB.QueryContext(*ctx, queryValidateUser)
	if err != nil {
		return status.Error(codes.Internal, "error with database. Details: "+err.Error())
	}

	defer rows.Close()
	if !rows.Next() {
		return status.Error(codes.NotFound, "user not found")
	}

	err = rows.Scan(&User.UserId, &User.Password)
	if err != nil {
		return status.Error(codes.Internal, "error scanning mysql row: "+err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(password))
	if err != nil {
		return status.Error(codes.InvalidArgument, "wrong password")
	}

	deleteQuery := fmt.Sprintf(`DELETE FROM users WHERE userId = '%s'`, User.UserId)
	_, err = mysql.DB.ExecContext(*ctx, deleteQuery)
	if err != nil {
		return status.Error(codes.Internal, "error with database. Details: "+err.Error())
	}

	return nil
}

func (r *mysqlResource) SendPasswordRecoveryEmail(ctx *context.Context, email string) (*string, error) {
	var User libModels.User

	queryValidateUser := fmt.Sprintf(`SELECT * FROM users WHERE email = '%s' LIMIT 1`, email)
	rows, err := mysql.DB.QueryContext(*ctx, queryValidateUser)
	if err != nil {
		return nil, status.Error(codes.Internal, "error with database. Details: "+err.Error())
	}

	defer rows.Close()
	if !rows.Next() {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	err = rows.Scan(&User.UserId, &User.Name, &User.Email, &User.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "error scanning mysql row: "+err.Error())
	}

	return &User.Name, nil
}

func (r *mysqlResource) GetUserProfile(ctx *context.Context, userId string) (*libModels.User, error) {
	var User libModels.User
	var settingsJSON []byte
	var settings libModels.UserSettings

	baseQuery := fmt.Sprintf(
		`SELECT 
    		u.name, 
    		u.email,
    		u.settings
		FROM
    		users u
		WHERE
    		u.userId = '%s'`,
		userId,
	)

	rows, err := mysql.DB.QueryContext(*ctx, baseQuery)
	if err != nil {
		return nil, status.Error(codes.Internal, "error with database. Details: "+err.Error())
	}

	defer rows.Close()
	if !rows.Next() {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	err = rows.Scan(&User.Name, &User.Email, &settingsJSON)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(settingsJSON, &settings)
	if err != nil {
		return nil, status.Error(codes.Internal, "error unmarshalling settings. Details: "+err.Error())
	}
	User.Settings = settings

	return &User, nil
}

func NewMysqlRepository(client *mysql.Client) IMySqlUser {
	return &mysqlResource{
		client: client,
	}
}

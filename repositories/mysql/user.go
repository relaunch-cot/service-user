package mysql

import (
	"context"
	"fmt"
	"github.com/relaunch-cot/lib-relaunch-cot/repositories/mysql"
)

type mysqlResource struct {
	client *mysql.Client
}

type IMySqlUser interface {
	CreateUser(ctx *context.Context, name, email, password string) error
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

func NewMysqlRepository(client *mysql.Client) IMySqlUser {
	return &mysqlResource{
		client: client,
	}
}

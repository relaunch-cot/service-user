package repositories

import (
	MysqlRepository "github.com/relaunch-cot/service-user/repositories/mysql"

	"github.com/relaunch-cot/lib-relaunch-cot/repositories/mysql"
)

type Repositories struct {
	Mysql MysqlRepository.IMySqlUser
}

func (r *Repositories) Inject(mysqlClient *mysql.Client) {
	r.Mysql = MysqlRepository.NewMysqlRepository(mysqlClient)
}

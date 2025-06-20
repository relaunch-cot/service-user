package resource

import (
	"github.com/relaunch-cot/service-user/handler"
	"github.com/relaunch-cot/service-user/repositories"
	"github.com/relaunch-cot/service-user/server"
)

var Repositories repositories.Repositories
var Handler handler.Handlers
var Server server.Servers

func Inject() {
	mysqlClient := OpenMysqlConn()

	Repositories.Inject(mysqlClient)
	Handler.Inject(&Repositories)
	Server.Inject(&Handler)
}

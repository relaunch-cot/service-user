package server

import (
	pb "github.com/relaunch-cot/lib-relaunch-cot/proto/user"
	"github.com/relaunch-cot/service-user/handler"
)

type Servers struct {
	User pb.UserServiceServer
}

func (s *Servers) Inject(handler *handler.Handlers) {
	s.User = NewUserServer(handler)
}

package methods

import (
	pb "github.com/relaunch-cot/lib-relaunch-cot/proto/user"
	"github.com/relaunch-cot/service-user/resource"
	"google.golang.org/grpc"
)

func RegisterGrpcServices(s *grpc.Server) {
	pb.RegisterUserServiceServer(s, resource.Server.User)
}

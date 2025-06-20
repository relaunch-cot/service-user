package main

import (
	"github.com/relaunch-cot/service-user/config"
	"github.com/relaunch-cot/service-user/resource"
	"github.com/relaunch-cot/service-user/server/methods"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	resource.Inject()

	lis, err := net.Listen("tcp", ":"+config.PORT)
	if err != nil {
		log.Fatalf("Failed to listen on %v: %v\n", config.PORT, err)
	}

	s := grpc.NewServer()

	methods.RegisterGrpcServices(s)

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}

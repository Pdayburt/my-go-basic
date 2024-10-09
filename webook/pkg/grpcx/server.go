package grpcx

import (
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	GrpcSvc *grpc.Server
	Addr    string
}

func (s *Server) Server() error {
	listen, err := net.Listen("tcp", ":8090")
	if err != nil {
		return err
	}
	//这边会阻塞
	return s.GrpcSvc.Serve(listen)

}

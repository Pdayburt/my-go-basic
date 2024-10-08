package grpc

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
	"time"
)

type Server struct {
	UnimplementedUserServiceServer
}

func (s *Server) GetById(ctx context.Context, req *GetByIdReq) (*GetByIdResp, error) {

	return &GetByIdResp{
		User: &User{
			Id:        123,
			Name:      "Jack",
			Avatar:    "sadadad.jpg",
			Attribute: map[string]string{"a": "b"},
		},
	}, nil
}

func TestServer(t *testing.T) {

	grpcSvc := grpc.NewServer()
	defer grpcSvc.GracefulStop()
	userSvc := &Server{}
	RegisterUserServiceServer(grpcSvc, userSvc)
	listener, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err)
	}
	err = grpcSvc.Serve(listener)
	if err != nil {
		t.Log(err)
	}

}

func TestClient(t *testing.T) {
	clientConn, err := grpc.Dial(":8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	client := NewUserServiceClient(clientConn)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	resp, err := client.GetById(ctx, &GetByIdReq{
		Id: 999,
	})
	assert.NoError(t, err)
	t.Log(resp.User)

}

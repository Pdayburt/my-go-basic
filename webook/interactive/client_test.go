package main

import (
	"context"
	intrv1 "example.com/mod/webook/api/proto/gen/intr/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"testing"
)

func TestGRPCClient(t *testing.T) {
	clientConn, err := grpc.Dial("localhost:8090",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer clientConn.Close()
	client := intrv1.NewInteractiveServiceClient(clientConn)
	//Get(ctx context.Context, in *GetReq, opts ...grpc.CallOption) (*GetResp, error)
	resp, err := client.Get(context.Background(), &intrv1.GetReq{
		Biz:   "test",
		BizId: 123,
		Uid:   345,
	})
	require.NoError(t, err)
	log.Println(resp.Intr)
}

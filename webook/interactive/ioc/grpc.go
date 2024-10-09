package ioc

import (
	"example.com/mod/webook/interactive/webookgrpc"
	"example.com/mod/webook/pkg/grpcx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitGRPCxServer(gIntrSvc *webookgrpc.InteractiveServiceServer) *grpcx.Server {
	type Config struct {
		Addr string `yaml:"addr"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		zap.L().Error("grpc 读取配置信息失败", zap.Error(err))
		panic(err)
	}
	grpcSvc := grpc.NewServer()
	gIntrSvc.Register(grpcSvc)
	return &grpcx.Server{Addr: cfg.Addr, GrpcSvc: grpcSvc}
}

package ioc

import (
	intrv1 "example.com/mod/webook/api/proto/gen/intr/v1"
	"example.com/mod/webook/interactive/service"
	"example.com/mod/webook/internal/client"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitIntrGRPCClient(intrService service.InteractiveService) intrv1.InteractiveServiceClient {
	type Config struct {
		Address string `yaml:"addr"`
		Secure  bool   `yaml:"secure"`
	}
	var cfg Config
	viper.UnmarshalKey("grpc.client.intr", &cfg)

	var opts []grpc.DialOption
	if cfg.Secure {

	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	dial, err := grpc.Dial(cfg.Address, opts...)
	if err != nil {
		panic(err)
	}
	remote := intrv1.NewInteractiveServiceClient(dial)
	LocalAdapter := client.NewInteractiveServiceAdapter(intrService)

	return client.NewGreyScaleInteractiveServiceClient(remote, LocalAdapter)

}

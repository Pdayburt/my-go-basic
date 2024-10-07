package config

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"testing"
)

func TestCfgRead(t *testing.T) {

	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	InitViperV1()
	fmt.Println(viper.AllSettings())
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg.Addrs)

}

// InitViperV1 使用结构体读取额配置文件
func InitViperV1() {
	viper.SetConfigFile("/Users/anatkh/Downloads/blockChain/golang/my-go-basic/webook/config/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

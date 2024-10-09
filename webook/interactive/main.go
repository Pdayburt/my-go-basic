package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func main() {
	InitViperV1()
	app := InitApp()
	consumers := app.consumer
	for _, consumer := range consumers {
		err := consumer.Start()
		if err != nil {
			panic(err)
		}
	}
	err := app.server.Server()

	if err != nil {
		panic(err)
	}

}

func InitViperV1() {
	viper.SetConfigFile("./config/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

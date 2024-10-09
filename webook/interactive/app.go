package main

import (
	"example.com/mod/webook/pkg/grpcx"
	"example.com/mod/webook/pkg/saramax"
)

type App struct {
	//所有需要在main函数里启动的，必须在这里
	//核心就是为了控制生命周期
	server   *grpcx.Server
	consumer []saramax.Consumer
}

package main

import (
	"example.com/mod/webook/internal/events"
	"github.com/gin-gonic/gin"
)

type App struct {
	server   *gin.Engine
	consumer []events.Consumer
}

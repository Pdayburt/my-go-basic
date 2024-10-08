package main

import (
	"example.com/mod/webook/internal/events"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

type App struct {
	server   *gin.Engine
	consumer []events.Consumer
	cron     *cron.Cron
}

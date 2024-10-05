package ioc

import (
	"example.com/mod/webook/internal/web"
	"example.com/mod/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func InitGin(mdls []gin.HandlerFunc, ud *web.UserHandler, ah *web.ArticleHandler) *gin.Engine {

	server := gin.Default()
	server.Use(mdls...)
	ud.RegisterRoutes(server)
	ah.RegisterRoutes(server)
	return server

}

func InitMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			//AllowOrigins: []string{"http://localhost:3000"},
			//AllowMethods:     []string{"POST", "GET", "OPTIONS"},
			ExposeHeaders:    []string{"x-jwt-token"},
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				if strings.Contains(origin, "http://localhost") ||
					strings.Contains(origin, "webook.com") {
					return true
				}
				return strings.Contains(origin, "company.com")
			},
			MaxAge: 12 * time.Hour,
		}),
		middleware.NewLoginJwtMiddlewareBuilder().
			IgnorePath("/users/login").
			IgnorePath("/users/signup").
			IgnorePath("/users/login_sms/code/send").
			IgnorePath("/users/login_sms").
			Build(),
	}
}

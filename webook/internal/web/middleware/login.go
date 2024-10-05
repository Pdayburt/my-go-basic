package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginMiddlewareBuilder struct {
	path []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}
func (lmb *LoginMiddlewareBuilder) IgnorePath(path string) *LoginMiddlewareBuilder {
	lmb.path = append(lmb.path, path)
	return lmb
}

func (lmb *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(context *gin.Context) {
		for _, v := range lmb.path {
			if context.Request.URL.Path == v {
				return
			}
		}
		session := sessions.Default(context)
		id := session.Get("userId")
		if id == nil {
			//没登陆
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

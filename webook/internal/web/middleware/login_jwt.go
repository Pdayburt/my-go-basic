package middleware

import (
	"example.com/mod/webook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

//JWT 登录校验

type LoginJwtMiddlewareBuilder struct {
	path []string
}

func NewLoginJwtMiddlewareBuilder() *LoginJwtMiddlewareBuilder {
	return &LoginJwtMiddlewareBuilder{}
}
func (lmb *LoginJwtMiddlewareBuilder) IgnorePath(path string) *LoginJwtMiddlewareBuilder {
	lmb.path = append(lmb.path, path)
	return lmb
}

func (lmb *LoginJwtMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(context *gin.Context) {
		for _, v := range lmb.path {
			if context.Request.URL.Path == v {
				return
			}
		}
		//使用jwt来校验
		tokenJWt := context.GetHeader("Authorization")
		if tokenJWt == "" {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		split := strings.Split(tokenJWt, " ")
		if len(split) != 2 {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := split[1]

		claims := &web.UserClaims{}
		/*token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte("30pzPuWsJCXJi5eryywAYltH5AS4GcOAA7aBwkDGpu0vGSqnVxjFLOmlLLNaWNsF"), nil
		})*/
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("30pzPuWsJCXJi5eryywAYltH5AS4GcOAA7aBwkDGpu0vGSqnVxjFLOmlLLNaWNsF"), nil
		})
		if err != nil {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != context.Request.UserAgent() {
			//	存在严重的安全问题
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//err不为nil token不为nil
		if !token.Valid || token == nil || claims.Uid == 0 {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//每十秒刷新一次 也就是每十秒重新生成jwt

		context.Set("claims", claims)
		context.Set("userId", claims.Uid)

	}
}

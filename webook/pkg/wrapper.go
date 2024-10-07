package pkg

import (
	"example.com/mod/webook/internal/web"
	"github.com/gin-gonic/gin"
	"net/http"
)

func WrapperBody[T any](fn func(ctx *gin.Context, params T) (web.Result, error)) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			return
		}
		//下半段的业务逻辑
		res, err := fn(ctx, req)
		if err != nil {
			//处理err 其实就是打印日志
		}
		ctx.JSON(http.StatusOK, res)

	}
}

package ratelimit

import (
	_ "embed"
	"example.com/mod/webook/pkg/ratelimit"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Builder struct {
	prefix string
	limit  ratelimit.Limiter
}

func NewBuilder(l ratelimit.Limiter) *Builder {
	return &Builder{
		prefix: "ip-limiter",
		limit:  l,
	}
}

func (b *Builder) Prefix(prefix string) *Builder {
	b.prefix = prefix
	return b
}

/*func (b *Builder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limited, err := b.limit(ctx)
		if err != nil {
			log.Println(err)
			// 这一步很有意思，就是如果这边出错了
			// 要怎么办？
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if limited {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}*/

func (b *Builder) Limit(ctx *gin.Context) (bool, error) {
	key := fmt.Sprintf("%s:%s", b.prefix, ctx.ClientIP())
	return b.limit.Limit(ctx, key)
}

package ratelimit

import "context"

type Limiter interface {
	// Limit limited 有没有出发限流 ，key就是限流对象
	//bool：是否限流
	//err 限流器本身有没有错误
	Limit(ctx context.Context, key string) (bool, error)
}

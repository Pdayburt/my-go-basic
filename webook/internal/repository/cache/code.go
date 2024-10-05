package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeSendTooMany   = errors.New("验证码发送太频繁")
	ErrCodeVerifyTooMany = errors.New("验证码验证太频繁")
	ErrUnknownForCode    = errors.New("未知错误")
)

//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type RedisCodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) CodeCache {
	return &RedisCodeCache{client: client}
}

func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	i, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int64()
	if err != nil {
		return err
	}
	switch i {
	case 0:
		return nil
	case -1:
		return ErrCodeSendTooMany
	default:
		return errors.New("系统错误")
	}
}

func (c *RedisCodeCache) Verify(ctx context.Context,
	biz, phone, inputCode string) (bool, error) {
	i, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, inputCode).Int64()
	if err != nil {
		return false, err
	}
	switch i {
	case 0:
		return true, nil
	case -1:
		return false, ErrCodeVerifyTooMany
	case -2:
		return false, nil
	default:
		return false, ErrUnknownForCode
	}

}

func (c *RedisCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

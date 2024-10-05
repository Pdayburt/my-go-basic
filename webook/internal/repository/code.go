package repository

import (
	"context"
	"example.com/mod/webook/internal/repository/cache"
)

var (
	ErrCodeSendTooMany = cache.ErrCodeSendTooMany
)

type CodeRepository interface {
	Store(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type CacheCodeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(c cache.CodeCache) CodeRepository {
	return &CacheCodeRepository{
		cache: c,
	}
}

func (cr *CacheCodeRepository) Store(ctx context.Context,
	biz, phone, code string) error {
	return cr.cache.Set(ctx, biz, phone, code)
}

func (c *CacheCodeRepository) Verify(ctx context.Context,
	biz, phone, inputCode string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, inputCode)

}

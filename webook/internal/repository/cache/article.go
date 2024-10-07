package cache

import (
	"context"
	"encoding/json"
	"example.com/mod/webook/internal/domain"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type ArticleCache interface {
	GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, uid int64, articles []domain.Article) error
	DeleteFirstPage(ctx context.Context, uid int64) error
	Set(ctx context.Context, uid int64, article domain.Article) error
}

type RedisArticleCache struct {
	client redis.Cmdable
}

func (r *RedisArticleCache) Set(ctx context.Context, uid int64, article domain.Article) error {
	return r.client.Set(ctx, r.key(uid), article, 1*time.Minute).Err()
}

func NewRedisArticleCache(client redis.Cmdable) ArticleCache {
	return &RedisArticleCache{client: client}
}

func (r *RedisArticleCache) GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error) {

	bs, err := r.client.Get(ctx, r.firstPageKey(uid)).Bytes()
	if err != nil {
		return nil, err
	}
	var articles []domain.Article
	err = json.Unmarshal(bs, &articles)
	return articles, err

}

func (r *RedisArticleCache) SetFirstPage(ctx context.Context, uid int64, articles []domain.Article) error {

	for i := 0; i < len(articles); i++ {
		articles[i].Content = articles[i].Abstract()
	}

	data, err := json.Marshal(articles)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.firstPageKey(uid), data,
		20*time.Second).Err()
}

func (r *RedisArticleCache) DeleteFirstPage(ctx context.Context, uid int64) error {
	return r.client.Del(ctx, r.firstPageKey(uid)).Err()
}

func (r *RedisArticleCache) firstPageKey(uid int64) string {
	return fmt.Sprintf("article:first_page:%d", uid)
}

func (r *RedisArticleCache) key(uid int64) string {
	return fmt.Sprintf("article:%d", uid)
}

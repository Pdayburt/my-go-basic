package cache

import (
	"context"
	"example.com/mod/webook/internal/domain"
	"fmt"
	"github.com/redis/go-redis/v9"

	"testing"
	"time"
)

func TestCache(t *testing.T) {

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:16379",
	})
	cache := NewRedisArticleCache(client)
	//ctx context.Context, uid int64, articles []domain.Article
	now := time.Now().UnixMilli()
	arts := []domain.Article{
		{
			Id:      99,
			Title:   "title",
			Content: "Content",
			//Author从用户来
			Author: domain.Author{
				Id:   98,
				Name: "jack",
			},
			Status: 0,
			Ctime:  time.UnixMilli(now),
			Utime:  time.UnixMilli(now),
		},
		{
			Id:      199,
			Title:   "1title",
			Content: "1Content",
			//Author从用户来
			Author: domain.Author{
				Id:   198,
				Name: "1jack",
			},
			Status: 10,
			Ctime:  time.UnixMilli(now),
			Utime:  time.UnixMilli(now),
		},
	}

	err := cache.SetFirstPage(context.Background(), int64(9628), arts)
	fmt.Println(err)

}

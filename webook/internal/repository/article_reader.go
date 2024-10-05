package repository

import (
	"context"
	"example.com/mod/webook/internal/domain"
)

type ArticleReaderRepository interface {
	//Save 有就更新 没有就新建
	Save(ctx context.Context, article domain.Article) (int64, error)
	//Update(ctx context.Context, article domain.Article) error
}

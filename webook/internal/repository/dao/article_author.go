package dao

import (
	"context"
	"gorm.io/gorm"
)

type ArticleAuthorDAO interface {
	Inert(ctx context.Context, article Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
}

func NewArticleAuthorDAO(db *gorm.DB) ArticleAuthorDAO {
	panic("implement me")
}

package dao

import (
	"context"
	"gorm.io/gorm"
)

type ArticleReaderDAO interface {
	Upsert(ctx context.Context, art Article) error
}

func NewArticleReaderDAO(db *gorm.DB) ArticleReaderDAO {
	panic("implement me")
}

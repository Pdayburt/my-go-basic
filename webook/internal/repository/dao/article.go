package dao

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDao interface {
	Insert(ctx context.Context, article Article) (int64, error)
	Update(ctx context.Context, article Article) error
	UpdateById(ctx context.Context, art Article) error
	Sync(ctx context.Context, article Article) (int64, error)
	Upsert(ctx context.Context, publishArticle PublicArticle) error
	Transaction(ctx context.Context, bizFun func(txDao ArticleDao) error) error
	SyncStatus(ctx *gin.Context, entity Article) error
}

type GORMArticleDAO struct {
	db *gorm.DB
}

func (G *GORMArticleDAO) SyncStatus(ctx *gin.Context, art Article) error {
	now := time.Now().UnixMilli()
	return G.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		res := tx.Model(&art).
			Where("id = ? AND author_id = ?", art.Id, art.AuthorId).
			Updates(map[string]any{
				"status": art.Status,
				"utime":  now,
			})
		if err := res.Error; err != nil {
			//数据库有问题
			return err
		}
		if res.RowsAffected != 1 {
			//id 或者 author_id出错
			return fmt.Errorf("有人在操作不是自己的文章,文章id：%d,author_id:%d", art.Id, art.AuthorId)
		}
		return tx.Model(&PublicArticle{}).
			Where("id = ?", art.Id).
			Updates(map[string]any{
				"status": art.Status,
				"utime":  now,
			}).Error

	})

}

// Sync 这是在dao层开启事务
func (G *GORMArticleDAO) Sync(ctx context.Context, article Article) (int64, error) {

	//线操作制作库（此时是表），后操作线上库（此时是表）
	//开启事务 这里采用了闭包形态--gorm帮忙我们管理类事务的生命周期 begin rollback commit由gorm自动管理
	var (
		id  = article.Id
		err error
	)
	err = G.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		txDAO := NewArticleDao(tx)
		if id > 0 {
			err = txDAO.UpdateById(ctx, article)
		} else {
			id, err = txDAO.Insert(ctx, article)
		}
		if err != nil {
			return err
		}

		//操作线上库
		return txDAO.Upsert(ctx, PublicArticle{article})

	})
	return id, err
}

func (G *GORMArticleDAO) Transaction(ctx context.Context, bizFun func(txDao ArticleDao) error) error {
	return G.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDAo := NewArticleDao(tx)
		return bizFun(txDAo)
	})
}

// Upsert update or insert
func (G *GORMArticleDAO) Upsert(ctx context.Context, publishArticle PublicArticle) error {

	now := time.Now().UnixMilli()
	publishArticle.Ctime = now
	publishArticle.Ctime = now

	err := G.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   publishArticle.Title,
			"content": publishArticle.Content,
			"status":  publishArticle.Status,
			"utime":   publishArticle.Utime,
		}),
	}).Create(&publishArticle).Error
	return err
}

func (G *GORMArticleDAO) UpdateById(ctx context.Context, art Article) error {

	now := time.Now().UnixMilli()
	tx := G.db.WithContext(ctx).Model(&art).
		Where("id = ? And author_id = ?", art.Id, art.AuthorId).
		Updates(map[string]any{
			"title":   art.Title,
			"content": art.Content,
			"utime":   now,
			"status":  art.Status,
		})
	if err := tx.Error; err != nil {
		return err
	}

	affected := tx.RowsAffected //更新行数
	if affected == 0 {
		return fmt.Errorf("更新失败，可能是创作者id非法:article id %d ,author_id: %d",
			art.Id, art.AuthorId)
	}
	return nil
}

func (G *GORMArticleDAO) Update(ctx context.Context, article Article) error {
	now := time.Now().UnixMilli()
	article.Utime = now
	err := G.db.WithContext(ctx).Model(&article).
		Where("id = ?", article.Id).
		Updates(map[string]any{"title": article.Title,
			"content": article.Content,
			"utime":   article.Utime,
		}).Error
	return err
}

func (G *GORMArticleDAO) Insert(ctx context.Context, article Article) (int64, error) {
	now := time.Now().UnixMilli()
	article.Ctime = now
	article.Utime = now
	err := G.db.WithContext(ctx).Create(&article).Error
	return article.AuthorId, err
}

// Article 如何设计索引？ 制作库的
type Article struct {
	Id      int64  `gorm:"primary_key;auto_increment"`
	Title   string `gorm:"type:varchar(1024)"`
	Content string `gorm:"type:BLOB"`

	AuthorId int64 `gorm:"index=aid_ctime"`
	Status   uint8 `gorm:"type:tinyint"`
	Ctime    int64 `gorm:"index=aid_ctime"`

	Utime int64
}

type PublicArticle struct {
	Article
}

func NewArticleDao(db *gorm.DB) ArticleDao {
	return &GORMArticleDAO{
		db: db,
	}
}

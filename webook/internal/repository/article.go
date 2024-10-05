package repository

import (
	"context"
	"example.com/mod/webook/internal/domain"
	"example.com/mod/webook/internal/repository/dao"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
	FindById(ctx context.Context, id int64) (domain.Article, error)
	Sync(ctx context.Context, article domain.Article) (int64, error)
	SyncStatus(ctx *gin.Context, article domain.Article) error
}

type CachedArticleRepository struct {
	// 操作单一的库
	dao dao.ArticleDao

	//在dao层操作两个数据源
	authorDAO dao.ArticleAuthorDAO
	readerDAO dao.ArticleReaderDAO
	db        *gorm.DB
}

func (c *CachedArticleRepository) SyncStatus(ctx *gin.Context, article domain.Article) error {

	return c.dao.SyncStatus(ctx, c.DomainToEntity(article))
}

func (c *CachedArticleRepository) Sync(ctx context.Context, article domain.Article) (int64, error) {

	return c.dao.Sync(ctx, c.DomainToEntity(article))

}

func (c *CachedArticleRepository) SyncV3(ctx context.Context, article domain.Article) (int64, error) {

	/*return c.dao.Transaction(ctx, func(txDao dao.ArticleDao) error {
		txDao.
	})*/
	panic("implement me")
}

// SyncV2 尝试在repository层面处理事务问题
func (c *CachedArticleRepository) SyncV2(ctx context.Context, article domain.Article) (int64, error) {
	//开启一个事务
	tx := c.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	defer tx.Rollback()
	//利用tx来构建DAO
	authorDAO := dao.NewArticleAuthorDAO(tx)
	readerDAO := dao.NewArticleReaderDAO(tx)
	var (
		id  = article.Id
		err error
	)
	entity := c.DomainToEntity(article)
	//先保存到制作库，再保存在线上库
	if id > 0 {
		id, err = authorDAO.Inert(ctx, entity)
	} else {
		id, err = authorDAO.Inert(ctx, entity)
	}
	if err != nil {
		tx.Rollback()
		return id, err
	}
	err = readerDAO.Upsert(ctx, entity)
	//执行成功直接提交
	if err != nil {
		tx.Commit()
	}
	return id, err

}

func (c *CachedArticleRepository) FindById(ctx context.Context, id int64) (domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CachedArticleRepository) Update(ctx context.Context, article domain.Article) error {
	return c.dao.UpdateById(ctx, dao.Article{
		Id:       article.Id,
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
		Status:   article.Status.ToUint8(),
	})
}

func (c *CachedArticleRepository) Create(ctx context.Context, article domain.Article) (int64, error) {
	return c.dao.Insert(ctx, dao.Article{
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
		Status:   article.Status.ToUint8(),
	})
}

func NewArticleRepository(dao dao.ArticleDao) ArticleRepository {
	return &CachedArticleRepository{
		dao: dao,
	}
}

func (c *CachedArticleRepository) DomainToEntity(article domain.Article) dao.Article {
	return dao.Article{
		Id:       article.Id,
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
		Status:   article.Status.ToUint8(),
	}
}

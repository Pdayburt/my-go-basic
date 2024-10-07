package repository

import (
	"context"
	"example.com/mod/webook/internal/domain"
	"example.com/mod/webook/internal/repository/cache"
	"example.com/mod/webook/internal/repository/dao"
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
	Sync(ctx context.Context, article domain.Article) (int64, error)
	SyncStatus(ctx context.Context, article domain.Article) error
	List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, uid int64, id int64) (domain.Article, error)
	GetPublishedById(ctx context.Context, uid int64, id int64) (domain.Article, error)
}

type CachedArticleRepository struct {
	// 操作单一的库
	articleDAO dao.ArticleDao
	userDao    dao.UserDao

	//在dao层操作两个数据源
	authorDAO dao.ArticleAuthorDAO
	readerDAO dao.ArticleReaderDAO
	db        *gorm.DB
	cache     cache.ArticleCache
}

func (c *CachedArticleRepository) GetPublishedById(ctx context.Context, uid int64, id int64) (domain.Article, error) {
	//读取线上数据，如果content已经放在oss上，则前端去读oss
	art, err := c.articleDAO.GetPubById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	//组装USer,适合单体
	user, err := c.userDao.FindById(ctx, art.AuthorId)
	res := c.EntityToDomain(art)
	res.Author.Name = user.NickName
	return res, err
}

func (c *CachedArticleRepository) GetById(ctx context.Context, uid int64, id int64) (domain.Article, error) {
	byId, err := c.articleDAO.GetById(ctx, uid, id)
	if err != nil {
		return domain.Article{}, err
	}
	return c.EntityToDomain(byId), nil
}

func (c *CachedArticleRepository) List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {

	//先访问缓存
	if offset == 0 && limit < 100 {
		data, err := c.cache.GetFirstPage(ctx, uid)
		if err == nil {
			return data, nil
		}
	}

	res, err := c.articleDAO.GetByAuthor(ctx, uid, offset, limit)
	if err != nil {
		return nil, err
	}

	artInDB := slice.Map[dao.Article, domain.Article](res, func(idx int, src dao.Article) domain.Article {
		return c.EntityToDomain(src)
	})

	//设置缓存
	err = c.cache.SetFirstPage(ctx, uid, artInDB)

	err = c.preCache(ctx, uid, artInDB)

	return artInDB, err

}

func (c *CachedArticleRepository) preCache(ctx context.Context, uid int64, arts []domain.Article) error {
	if len(arts) > 0 {
		return c.cache.Set(ctx, uid, arts[0])
	}
	return fmt.Errorf("缓存的内容不能为空")
}

func (c *CachedArticleRepository) SyncStatus(ctx context.Context, article domain.Article) error {

	return c.articleDAO.SyncStatus(ctx, c.DomainToEntity(article))
}

func (c *CachedArticleRepository) Sync(ctx context.Context, article domain.Article) (int64, error) {

	//如果数据更新 则清空缓存
	defer func() {
		err := c.cache.DeleteFirstPage(ctx, article.Author.Id)
		if err != nil {
			zap.L().Error("更新article时删除缓存失败", zap.Error(err))
		}
	}()
	return c.articleDAO.Sync(ctx, c.DomainToEntity(article))

}

func (c *CachedArticleRepository) SyncV3(ctx context.Context, article domain.Article) (int64, error) {

	/*return c.dao.Transaction(ctx, func(txDao dao.ArticleDao) error {
		txDao.
	})*/
	panic("implement me")
}

// SyncV2 尝试在repository层面处理事务问题
func (c *CachedArticleRepository) SyncV2(ctx context.Context, article domain.Article) (int64, error) {
	//如果数据更新 则清空缓存
	defer func() {
		err := c.cache.DeleteFirstPage(ctx, article.Author.Id)
		if err != nil {
			zap.L().Error("更新article时删除缓存失败", zap.Error(err))
		}
	}()
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

func (c *CachedArticleRepository) Update(ctx context.Context, article domain.Article) error {

	//如果数据更新 则清空缓存
	defer func() {
		err := c.cache.DeleteFirstPage(ctx, article.Author.Id)
		if err != nil {
			zap.L().Error("更新article时删除缓存失败", zap.Error(err))
		}
	}()

	return c.articleDAO.UpdateById(ctx, dao.Article{
		Id:       article.Id,
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
		Status:   article.Status.ToUint8(),
	})
}

func (c *CachedArticleRepository) Create(ctx context.Context, article domain.Article) (int64, error) {

	//如果数据更新 则清空缓存
	defer func() {
		err := c.cache.DeleteFirstPage(ctx, article.Author.Id)
		if err != nil {
			zap.L().Error("更新article时删除缓存失败", zap.Error(err))
		}
	}()

	return c.articleDAO.Insert(ctx, dao.Article{
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
		Status:   article.Status.ToUint8(),
	})
}

func NewArticleRepository(dao dao.ArticleDao) ArticleRepository {
	return &CachedArticleRepository{
		articleDAO: dao,
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
func (c *CachedArticleRepository) EntityToDomain(article dao.Article) domain.Article {

	return domain.Article{
		Id:      article.Id,
		Title:   article.Title,
		Content: article.Content,
		Author:  domain.Author{Id: article.AuthorId},
		Ctime:   time.UnixMilli(article.Ctime),
		Utime:   time.UnixMilli(article.Utime),
	}
}

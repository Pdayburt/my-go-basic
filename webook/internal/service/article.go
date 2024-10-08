package service

import (
	"context"
	"example.com/mod/webook/internal/domain"
	"example.com/mod/webook/internal/events/article"
	"example.com/mod/webook/internal/repository"
	"go.uber.org/zap"
)

//go:generate mockgen -source=./article.go -package=svcmocks -destination=svcmocks/article.mock.go ArticleService
type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error)
	Withdraw(ctx context.Context, article domain.Article) error
	List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, uid int64, id int64) (domain.Article, error)
	GetPublishedById(ctx context.Context, uid int64, id int64) (domain.Article, error)
	ListPub(ctx context.Context, offset int, limit int) ([]domain.Article, error)
}

type articleService struct {
	repo repository.ArticleRepository
	//
	author   repository.ArticleAuthorRepository
	reader   repository.ArticleReaderRepository
	producer article.Producer
}

func (as *articleService) ListPub(ctx context.Context, offset int, limit int) ([]domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

/*func NewArticleService(author repository.ArticleAuthorRepository, reader repository.ArticleReaderRepository) ArticleService {
	return &articleService{
		author: author,
		reader: reader,
	}
}*/

func NewArticleService(repo repository.ArticleRepository, producer article.Producer) ArticleService {
	return &articleService{
		repo:     repo,
		producer: producer,
	}
}

func (as *articleService) GetPublishedById(ctx context.Context, uid int64, id int64) (domain.Article, error) {
	art, err := as.repo.GetPublishedById(ctx, uid, id)
	if err != nil {
		go func() {
			er := as.producer.ProducerReadEvent(ctx, article.ReadEvent{
				Uid: uid,
				Aid: id,
			})
			if er != nil {
				zap.L().Error("发送消息至kafka失败", zap.Error(err))
			}
		}()
	}

	return art, err

}

func (as *articleService) GetById(ctx context.Context, uid int64, id int64) (domain.Article, error) {
	return as.repo.GetById(ctx, uid, id)
}

func (as *articleService) List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	return as.repo.List(ctx, uid, offset, limit)
}

func (as *articleService) Withdraw(ctx context.Context, article domain.Article) error {
	article.Status = domain.ArticleStatusPrivate
	return as.repo.SyncStatus(ctx, article)
}

func (as *articleService) Publish(ctx context.Context, article domain.Article) (int64, error) {
	//先保存到制作库
	//再保存到线上库
	/*
		在service层调用不同的数据库


		var (
			id  = article.Id
			err error
		)

		if id > 0 {
			err = as.author.Update(ctx, article)
		} else {
			id, err = as.author.Create(ctx, article)
		}
		if err != nil {
			return 0, err
		}
		//两个库的id需要保持一致
		article.Id = id

		//引入重试机制
		for i := 0; i < 3; i++ {
			id, err = as.reader.Save(ctx, article)
			if err == nil {
				break
			}

			zap.L().Error("保存至线上库失败",
				zap.Int64("文章的id ：", id), zap.Error(err))

		}
		if err != nil {
			zap.L().Error("部分失败：数据保存到线上库重试也失败",
				zap.Int64("文章的id ：", id), zap.Error(err))
		}
		//接入告警心痛 手工处理
		return id, err*/
	article.Status = domain.ArticleStatusPublished
	//以下是在dao层处理两个数据库
	return as.repo.Sync(ctx, article)
}

func (as *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {

	article.Status = domain.ArticleStatusPublished
	if article.Id > 0 {
		err := as.repo.Update(ctx, article)
		return article.Id, err
	}
	return as.repo.Create(ctx, article)

	/*	if article.Id > 0 {
			//进行更新操作
			err := as.author.Update(ctx, article)
			return article.Id, err
		}
		//进行保存工作
		return as.author.Create(ctx, article)
	}*/

	/*func (as *articleService) Update(ctx context.Context, article domain.Article) error {
	articleInDB, err := as.repo.FindById(ctx, article.Id)
	if err != nil {
		return fmt.Errorf("根据Id查询文章失败 : %w", err)
	}
	if article.Id != articleInDB.Author.Id {
		return fmt.Errorf("禁止更新别人的文章 ")
	}
	return as.repo.Update(ctx, article)
	return 0, nil*/
}

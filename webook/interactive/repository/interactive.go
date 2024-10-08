package repository

import (
	"context"
	"example.com/mod/webook/interactive/domain"
	"example.com/mod/webook/interactive/repository/cache"
	"example.com/mod/webook/interactive/repository/dao"
	"go.uber.org/zap"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, id int64) error
	BatchIncrReadCnt(ctx context.Context, biz []string, id []int64) error
	IncrLike(ctx context.Context, biz string, id int64, uid int64) error
	DecrLike(ctx context.Context, biz string, id int64, uid int64) error
	AddCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error
	Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	Get(ctx context.Context, biz string, id int64) (domain.Interactive, error)
}

type interactiveRepository struct {
	dao   dao.InteractiveDao
	cache cache.InteractiveCache
}

func (i *interactiveRepository) BatchIncrReadCnt(ctx context.Context, biz []string, id []int64) error {

	return i.dao.BatchIncrReadCnt(ctx, biz, id)
}

func (i *interactiveRepository) Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := i.dao.GetCollectionInfo(ctx, biz, id, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (i *interactiveRepository) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {

	//先查缓存
	interactive, err := i.cache.Get(ctx, biz, bizId)
	if err == nil {
		return interactive, err
	}

	//再查询数据库
	daoIntr, err := i.dao.Get(ctx, biz, bizId)
	if err != nil {
		return domain.Interactive{}, err
	}
	intr := i.toDomain(daoIntr)
	go func() {
		er := i.cache.Set(ctx, biz, bizId, intr)
		//记录日志
		if er != nil {
			zap.L().Error("回写缓存失败",
				zap.String("biz", biz),
				zap.Error(er))
		}
	}()
	return intr, nil
}

func (i *interactiveRepository) Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error) {

	_, err := i.dao.GetLikeInfo(ctx, biz, id, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (i *interactiveRepository) AddCollectionItem(ctx context.Context, biz string, id int64, cid int64, uid int64) error {
	err := i.dao.InsertCollectionBiz(ctx, dao.UserCollectionBiz{
		Cid:   cid,
		BizId: id,
		Biz:   biz,
		Uid:   uid,
	})
	if err != nil {
		return err
	}
	//收藏个数
	return i.cache.IncrCollectionCntIfPresent(ctx, biz, id)
}

func (i *interactiveRepository) IncrLike(ctx context.Context, biz string, id int64, uid int64) error {
	//县插入点赞 再更新点赞数 再更新缓存
	err := i.dao.InsertLikeInfo(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	return i.cache.IncrReadCntIfPresent(ctx, biz, id)
}

func (i *interactiveRepository) DecrLike(ctx context.Context, biz string, id int64, uid int64) error {

	err := i.dao.DeleteLikeInfo(ctx, biz, id, uid)
	if err != nil {
		return err
	}
	return i.cache.DecrReadCntIfPresent(ctx, biz, id)
}

func (i *interactiveRepository) IncrReadCnt(ctx context.Context, biz string, id int64) error {

	//需要考虑缓存方案
	err := i.dao.IncrReadCnt(ctx, biz, id)
	if err != nil {
		return err
	}

	return i.cache.IncrReadCntIfPresent(ctx, biz, id)
}

func NewInteractiveRepository(dao dao.InteractiveDao, cache cache.InteractiveCache) InteractiveRepository {
	return &interactiveRepository{dao: dao, cache: cache}
}

func (c *interactiveRepository) toDomain(intr dao.Interactive) domain.Interactive {
	return domain.Interactive{
		Biz:        intr.Biz,
		BizId:      intr.BizId,
		LikeCnt:    intr.LikeCnt,
		CollectCnt: intr.CollectCnt,
		ReadCnt:    intr.ReadCnt,
	}
}

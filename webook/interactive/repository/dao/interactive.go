package dao

import (
	"context"
	"go.uber.org/zap"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

var ErrRecordNotFound = gorm.ErrRecordNotFound

type InteractiveDao interface {
	IncrReadCnt(ctx context.Context, biz string, id int64) error
	InsertLikeInfo(ctx context.Context, biz string, id int64, uid int64) error
	DeleteLikeInfo(ctx context.Context, biz string, id int64, uid int64) error
	InsertCollectionBiz(ctx context.Context, biz UserCollectionBiz) error
	GetLikeInfo(ctx context.Context, biz string, id int64, uid int64) (UserLikeBiz, error)
	GetCollectionInfo(ctx context.Context, biz string, id int64, uid int64) (UserCollectionBiz, error)
	Get(ctx context.Context, biz string, id int64) (Interactive, error)
	BatchIncrReadCnt(ctx context.Context, biz []string, id []int64) error
	GetByIds(ctx context.Context, biz string, ids []int64) ([]Interactive, error)
}

type interactiveDao struct {
	db *gorm.DB
}

func (i *interactiveDao) GetByIds(ctx context.Context, biz string, ids []int64) ([]Interactive, error) {
	var res []Interactive
	err := i.db.WithContext(ctx).Where("biz = ? AND id IN ?", biz, ids).Find(&res).Error
	return res, err
}

func NewInteractiveDao(db *gorm.DB) InteractiveDao {
	return &interactiveDao{db: db}
}

func (i *interactiveDao) BatchIncrReadCnt(ctx context.Context, biz []string, id []int64) error {

	i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDAO := NewInteractiveDao(tx)

		for j := 0; j < len(biz); j++ {
			err := txDAO.IncrReadCnt(ctx, biz[j], id[j])
			if err != nil {
				//记录日志即可 因为这个数据不是很重要
				zap.L().Error("kafka中的数据插入失败", zap.Int64("id:", id[j]))
			}
		}

		return nil
	})
	return nil
}

func (i *interactiveDao) Get(ctx context.Context, biz string, id int64) (Interactive, error) {
	var intr Interactive
	err := i.db.WithContext(ctx).
		Where("biz = ? and id = ?", biz, id).
		Find(&intr).Error
	return intr, err
}

func (i *interactiveDao) GetCollectionInfo(ctx context.Context, biz string, id int64, uid int64) (UserCollectionBiz, error) {
	var res UserCollectionBiz
	err := i.db.WithContext(ctx).
		Where("biz=? AND biz_id = ? AND uid = ?", biz, id, uid).
		First(&res).Error
	return res, err

}

func (i *interactiveDao) GetLikeInfo(ctx context.Context, biz string, id int64, uid int64) (UserLikeBiz, error) {

	var user UserLikeBiz
	err := i.db.WithContext(ctx).
		Where("biz=? AND biz_id = ? AND uid = ? AND status = ?",
			biz, id, uid, 1).
		Find(&user).Error
	return user, err
}

func (i *interactiveDao) InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error {
	now := time.Now().UnixMilli()
	cb.Utime = now
	cb.Ctime = now
	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := i.db.WithContext(ctx).Create(&cb).Error
		if err != nil {
			return err
		}
		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"collect_cnt": gorm.Expr("`collect_cnt`+1"),
				"utime":       now,
			}),
		}).Create(&Interactive{
			CollectCnt: 1,
			Ctime:      now,
			Utime:      now,
			Biz:        cb.Biz,
			BizId:      cb.BizId,
		}).Error
	})

}

func (i *interactiveDao) DeleteLikeInfo(ctx context.Context, biz string, id int64, uid int64) error {

	now := time.Now().UnixMilli()

	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		//两个操作 软删除点赞记录和减少点赞数量

		err := tx.Model(&UserLikeBiz{}).
			Where("biz = ? and biz_id = ? and uid =?", id, biz, uid).
			Updates(map[string]any{
				"utime":  now,
				"status": 0,
			}).Error

		if err != nil {
			return err
		}

		return tx.Model(&Interactive{}).
			Where("biz = ? and bit_id = ?", biz, id).
			Updates(map[string]any{
				"utime":    now,
				"like_cnt": gorm.Expr("like_cnt - 1"),
			}).Error

	})
}

func (i *interactiveDao) InsertLikeInfo(ctx context.Context, biz string, id int64, uid int64) error {
	//
	now := time.Now().UnixMilli()
	err := i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"status": 1,
				"utime":  now,
			})},
		).
			Create(&UserLikeBiz{
				Biz:    biz,
				BizId:  id,
				Uid:    uid,
				Ctime:  now,
				Utime:  now,
				Status: 1,
			}).Error
		if err != nil {
			return err
		}
		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"like_cnt": gorm.Expr("like_cnt + ?", 1),
				"utime":    now,
			}),
		}).Create(&Interactive{
			LikeCnt: 1,
			Ctime:   now,
			Utime:   now,
			Biz:     biz,
			BizId:   id,
		}).Error

	})
	return err
}

func (i *interactiveDao) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {

	//dao要如何实现
	//mysql 支持原子操作 a = a+1
	now := time.Now().UnixMilli()
	return i.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"read_cnt": gorm.Expr("read_cnt +1"),
			"utime":    time.Now().UnixMilli(),
		}),
	}).Create(&Interactive{
		BizId:   bizId,
		Biz:     biz,
		ReadCnt: 1,
		Ctime:   now,
		Utime:   now,
	}).Error

	/*return i.db.WithContext(ctx).
	Where("biz_id = ? and biz =?", bizId, biz).
	Updates(map[string]any{
		"read_cnt": gorm.Expr("read_cnt +1"),
		"utime":    time.Now().UnixMilli(),
	}).Error
	*/
}

type UserLikeBiz struct {
	Id    int64  `gorm:"primaryKey,autoIncrement:false"`
	Biz   string `gorm:"type:varchar(128);uniqueIndex:idx_uid_biz_id_biz"`
	BizId int64  `gorm:"uniqueIndex:idx_uid_biz_id_biz"`
	//谁的操作
	Uid   int64 `gorm:"uniqueIndex:idx_uid_biz_id_biz"`
	Ctime int64
	Utime int64
	//软删除
	Status uint8
}

type Interactive struct {
	Id         int64  `gorm:"primaryKey,autoIncrement"`
	BizId      int64  `gorm:"uniqueIndex:biz_type_id"`
	Biz        string `gorm:"type:varchar(128);uniqueIndex:biz_type_id"`
	ReadCnt    int64
	LikeCnt    int64
	CollectCnt int64
	Ctime      int64
	Utime      int64
}

// Error 1170 (42000): BLOB/TEXT column 'biz' used in key specification without a key length
// Collection 收藏夹
type Collection struct {
	Id   int64  `gorm:"primaryKey,autoIncrement"`
	Name string `gorm:"type=varchar(1024)"`
	Uid  int64  `gorm:""`

	Ctime int64
	Utime int64
}

// UserCollectionBiz 收藏的东西
type UserCollectionBiz struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 收藏夹 ID
	// 作为关联关系中的外键，我们这里需要索引
	Cid   int64  `gorm:"index"`
	BizId int64  `gorm:"uniqueIndex:biz_type_id_uid"`
	Biz   string `gorm:"type:varchar(128);uniqueIndex:biz_type_id_uid"`
	// 这算是一个冗余，因为正常来说，
	// 只需要在 Collection 中维持住 Uid 就可以
	Uid   int64 `gorm:"uniqueIndex:biz_type_id_uid"`
	Ctime int64
	Utime int64
}

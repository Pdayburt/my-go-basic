package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDao interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
}

type GORMUserDAO struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDao {
	return &GORMUserDAO{
		db: db,
	}
}

// 对应与数据库表结构一一对应
type User struct {
	Id       int64          `gorm:"primaryKey,AUTO_INCREMENT"`
	Email    sql.NullString `gorm:"unique"`
	Password string
	Ctime    int64
	Utime    int64
	Phone    sql.NullString `gorm:"unique"`
}

func (ud *GORMUserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := ud.db.WithContext(ctx).Create(&u).Error
	var mysqErr *mysql.MySQLError
	if errors.As(err, &mysqErr) {
		const uniqueConflictErr uint16 = 1062
		if mysqErr.Number == uniqueConflictErr {
			//邮箱冲突
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (ud *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := ud.db.WithContext(ctx).First(&user, "email = ?", email).Error
	return user, err

}

func (ud *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var user User
	err := ud.db.WithContext(ctx).First(&user, "phone = ?", phone).Error
	return user, err

}

func (ud *GORMUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var user User
	err := ud.db.WithContext(ctx).Find(&user, "id = ?", id).Error
	return user, err
}

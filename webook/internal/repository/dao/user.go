package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	//这里存的是毫秒数，也可以存秒数，或者纳秒数，看个人喜好
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			//	邮箱冲突
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	//err := dao.db.WithContext(ctx).First(&u,"email = ?",email).Error
	return u, err
}

// User 直接对应于数据库表结构
// 有些人叫做 entity,有些人叫做 model 也有人叫做 PO(persistent object)
type User struct {
	// primaryKey 指定列为主键
	// autoIncrement 指定列为自动增长
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 全部用户唯一
	// unique 指定列为唯一，不能有相同的
	Email    string `gorm:"unique"`
	Password string

	// 	往这里面加

	//创建时间，毫秒数
	Ctime int64
	//更新时间，毫秒数
	Utime int64
}

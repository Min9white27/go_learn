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
	ErrUserDuplicate = errors.New("邮箱或者手机号码冲突")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	// FindById FindByIdV1(ctx context.Context, uid interface{}) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	//UpdateByUid(ctx context.Context, entity User) error
}

type GORMUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}

func (dao *GORMUserDAO) Insert(ctx context.Context, u User) error {
	//这里存的是毫秒数，也可以存秒数，或者纳秒数，看个人喜好
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			//	邮箱冲突 or 手机号码冲突
			return ErrUserDuplicate
		}
	}
	return err
}

func (dao *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	//err := dao.db.WithContext(ctx).First(&u,"email = ?",email).Error
	return u, err
}

func (dao *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	//err := dao.db.WithContext(ctx).First(&u,"email = ?",email).Error
	return u, err
}

//func (dao *GORMUserDAO) FindByIdV1(ctx context.Context, uid interface{}) (User, error) {
//	var u User
//	err := dao.db.WithContext(ctx).Where("id = ?", uid).First(&u).Error
//	return u, err
//}

func (dao *GORMUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("'id' = ?", id).First(&u).Error
	return u, err
}

//func (dao *GORMUserDAO) UpdateByUid(ctx context.Context, entity User) error {
//	return dao.db.WithContext(ctx).Model(&entity).Where("id = ?", entity.Id).
//		Updates(map[string]any{
//			"utime":            time.Now().UnixMilli(),
//			"nickname":         entity.Nickname,
//			"birthday":         entity.Birthday,
//			"personal_profile": entity.PersonalProfile,
//		}).Error
//}

// User 直接对应于数据库表结构
// 有些人叫做 entity,有些人叫做 model 也有人叫做 PO(persistent object)
type User struct {
	// primaryKey 指定列为主键
	// autoIncrement 指定列为自动增长
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 全部用户唯一
	// unique 指定列为唯一，不能有相同的
	Email    sql.NullString `gorm:"unique"`
	Password string

	// 唯一索引允许有多个空值，但是不能有多个空字符串 ""
	Phone sql.NullString `gorm:"unique"`
	// 这种理论上是可行的，但有一个很大的问题就是要解引用，解引用的问题就是要判空
	//phone *string

	// 	往这里面加
	//`gorm:"type=varchar(128)"`
	//Nickname string
	//Birthday int64
	// 指定是 varchar 这个类型的，并且长度是 1024
	//`gorm:"type=varchar(4096)"`
	//PersonalProfile string

	//创建时间，毫秒数
	Ctime int64
	//更新时间，毫秒数
	Utime int64
}

package repository

import (
	"context"
	"gitee.com/geektime-geekbang_admin/geektime-basic-go/webook/internal/domain"
	"gitee.com/geektime-geekbang_admin/geektime-basic-go/webook/internal/repository/dao"
	"time"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	// SELECT * FROM `users` WHERE `email`=?
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) FindById(ctx context.Context, uid interface{}) (domain.User, error) {
	//先从 cache 缓存里面找
	//再从 dao 里面找
	//找到了回写 cache 缓存
	u, err := r.dao.FindById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	return r.toDomain(u), nil
}

func (r *UserRepository) UpdateNonZeroFields(ctx context.Context, user domain.User) error {
	return r.dao.UpdateByUid(ctx, r.toEntity(user))
}

func (r *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:              u.Id,
		Email:           u.Email,
		Password:        u.Password,
		Nickname:        u.Nickname,
		Birthday:        time.UnixMilli(u.Birthday),
		PersonalProfile: u.PersonalProfile,
	}
}

func (r *UserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		Id:              u.Id,
		Email:           u.Email,
		Password:        u.Password,
		Nickname:        u.Nickname,
		Birthday:        u.Birthday.UnixMilli(),
		PersonalProfile: u.PersonalProfile,
	}
}
